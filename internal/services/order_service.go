package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/types"
	"github.com/moonrill/rumahpc-api/utils"
	"github.com/xendit/xendit-go/v6/invoice"
	"gorm.io/gorm"
)

var (
	ErrOrderProductNotFound   = errors.New("order product not found")
	ErrOrderAddressNotFound   = errors.New("order address not found")
	ErrOrderShipping          = errors.New("error shipping order")
	ErrOrderCreate            = errors.New("error creating order")
	ErrInvalidShippingOptions = errors.New("invalid shipping options")
	ErrInvalidQuantity        = errors.New("invalid quantity")
	ErrOrderCartItemNotFound  = errors.New("one or more cart items not found")
	ErrCartItemRemoval        = errors.New("failed to remove cart items")
)

func CreateBuyNowOrder(request *types.BuyNowRequest, userID string) (*invoice.Invoice, error) {
	// Get user address
	var address models.Address
	err := config.DB.First(&address, "id = ?", request.AddressID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrOrderAddressNotFound
		}
		return nil, err
	}

	var user models.User
	err = config.DB.First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, utils.ErrNotFound
	}

	var product models.Product
	err = config.DB.Preload("Category").
		First(&product, "id = ?", request.ProductID).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrOrderProductNotFound
		}
		return nil, err
	}

	// Check if product is active & stock > 0
	if product.Status != models.ProductStatusActive || product.Stock <= 0 {
		return nil, utils.ErrProductUnavailable
	}

	// Calculate Total Price
	totalPrice := product.Price * request.Quantity

	// Create Order
	orderID := uuid.New().String()
	order := &models.Order{
		ID:             orderID,
		UserID:         userID,
		AddressID:      request.AddressID,
		Status:         models.OrderStatusWaiting,
		ShippingPrice:  &request.ShippingPrice,
		TotalPrice:     totalPrice + request.ShippingPrice,
		CourierCompany: &request.CourierCompany,
		CourierType:    &request.CourierType,
		User:           &user,
	}

	err = config.DB.Create(order).Error
	if err != nil {
		return nil, ErrOrderCreate
	}

	// Create Order Items
	orderItem := &models.OrderItem{
		OrderID:   orderID,
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
		SubTotal:  product.Price * request.Quantity,
	}

	err = config.DB.Create(orderItem).Error
	if err != nil {
		return nil, ErrOrderCreate
	}

	// Create Xendit Invoice
	inv, err := CreateBuyNowXenditInvoice(order, orderItem)

	if err != nil {
		return nil, ErrOrderCreate
	}

	return inv, nil
}

func CreateCartCheckoutOrder(request *types.CheckoutCartRequest, user *models.User) (*invoice.Invoice, error) {
	// Get user address
	var address models.Address
	err := config.DB.First(&address, "id = ?", request.AddressID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrOrderAddressNotFound
		}
		return nil, err
	}

	// Fetch cart items with their associated products
	var cartItems []models.CartItem
	err = config.DB.Preload("Product").Where("id IN (?)", request.CartItems).Find(&cartItems).Error
	if err != nil {
		return nil, ErrOrderCartItemNotFound
	}

	// Validate cart items belong to the user and exist
	if len(cartItems) != len(request.CartItems) {
		return nil, ErrOrderCartItemNotFound
	}

	// Group cart items by merchant
	merchantCartItems := make(map[string][]models.CartItem)
	for _, cartItem := range cartItems {
		if cartItem.Product == nil {
			return nil, ErrOrderProductNotFound
		}
		merchantCartItems[cartItem.Product.MerchantID] = append(
			merchantCartItems[cartItem.Product.MerchantID],
			cartItem,
		)
	}

	// Validate Shipping Options match merchant groups
	if len(merchantCartItems) != len(request.ShippingOptions) {
		return nil, ErrInvalidShippingOptions
	}

	var orders []*models.Order
	var orderItems []*models.OrderItem
	var totalOrderPrice int

	// Create orders for each merchant group
	for merchantID, cartItemGroup := range merchantCartItems {
		// Find matching shipping option
		var shippingOption *types.ShippingOption
		for _, opt := range request.ShippingOptions {
			if opt.MerchantID == merchantID {
				shippingOption = &opt
				break
			}
		}

		if shippingOption == nil {
			return nil, ErrInvalidShippingOptions
		}

		orderID := uuid.New().String()

		var merchantOrderTotal int

		for _, cartItem := range cartItemGroup {
			// Validate quantity
			if cartItem.Quantity < 1 {
				return nil, ErrInvalidQuantity
			}

			// Calculate item subtotal
			itemSubTotal := cartItem.Product.Price * cartItem.Quantity
			merchantOrderTotal += itemSubTotal

			// Create order item
			orderItem := models.OrderItem{
				OrderID:   orderID,
				ProductID: cartItem.ProductID,
				Quantity:  cartItem.Quantity,
				SubTotal:  itemSubTotal,
			}
			orderItems = append(orderItems, &orderItem)
		}

		// Create order
		order := models.Order{
			ID:             orderID,
			UserID:         user.ID,
			AddressID:      request.AddressID,
			Status:         models.OrderStatusWaiting,
			ShippingPrice:  &shippingOption.Price,
			TotalPrice:     merchantOrderTotal + shippingOption.Price,
			CourierCompany: &shippingOption.CourierCompany,
			CourierType:    &shippingOption.CourierType,
			User:           user,
		}

		totalOrderPrice += order.TotalPrice
		orders = append(orders, &order)
	}

	// Begin database transaction
	tx := config.DB.Begin()
	if tx.Error != nil {
		return nil, ErrOrderCreate
	}

	// Create orders
	for _, order := range orders {
		if err := tx.Create(order).Error; err != nil {
			tx.Rollback()
			return nil, ErrOrderCreate
		}
	}

	// Create order items
	for _, orderItem := range orderItems {
		if err := tx.Create(orderItem).Error; err != nil {
			tx.Rollback()
			return nil, ErrOrderCreate
		}
	}

	// Remove cart items after successful order creation
	if err := tx.Where("id IN (?)", request.CartItems).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		return nil, ErrCartItemRemoval
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, ErrOrderCreate
	}

	// Create Xendit Invoice
	inv, err := CreateCartCheckoutXenditInvoice(orders, orderItems)
	if err != nil {
		return nil, ErrOrderCreate
	}

	return inv, nil
}
