package services

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/templates"
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
		TotalPrice:     totalPrice,
		CourierCompany: &request.CourierCompany,
		CourierType:    &request.CourierType,
		User:           &user,
		MerchantID:     product.MerchantID,
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
			TotalPrice:     merchantOrderTotal,
			CourierCompany: &shippingOption.CourierCompany,
			CourierType:    &shippingOption.CourierType,
			User:           user,
			MerchantID:     merchantID,
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

func GetOrders(userID string, page, limit int, status, shippingStatus string) ([]*models.Order, int64, error) {
	var orders []*models.Order
	var totalCount int64

	offset := (page - 1) * limit

	query := config.DB.Model(&models.Order{}).Where("user_id = ?", userID)

	// Apply filters for status and shippingStatus if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if shippingStatus != "" {
		query = query.Where("shipping_status = ?", shippingStatus)
	}

	// Count total records with applied filters
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Fetch orders with applied filters
	result := query.
		Preload("Merchant").
		Preload("OrderItems.Product").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&orders)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return orders, totalCount, nil
}

func GetOrderById(userId, orderId string) (*models.Order, error) {
	var order models.Order
	err := config.DB.
		Preload("Merchant").
		Preload("OrderItems.Product").
		Preload("Address").
		Preload("Payment").
		First(&order, "user_id = ? AND id = ?", userId, orderId).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	return &order, nil
}

func CompleteOrder(orderID string, userID string) error {
	order, err := GetOrderById(userID, orderID)
	if err != nil {
		return err
	}

	if order.Status != models.OrderStatusDelivered {
		return utils.ErrBadRequest
	}

	order.Status = models.OrderStatusCompleted

	// Create Payout
	err = CreatePayoutOrder(order)
	if err != nil {
		return err
	}

	return config.DB.Save(order).Error
}

func CancelOrder(orderID string, userID string) error {
	order, err := GetOrderById(userID, orderID)
	if err != nil {
		return err
	}

	if order.Status != models.OrderStatusProcessing && order.Status != models.OrderStatusShipped {
		return utils.ErrBadRequest
	}

	order.Status = models.OrderStatusCancelled
	order.ShippingStatus = "cancelled"

	if err := CancelShipping(*order.ShippingID); err != nil {
		return err
	}

	// TODO: Refund Payment

	return config.DB.Save(order).Error
}

func GenerateInvoicePdf(orderID string, user *models.User) ([]byte, error) {
	// Get order details
	order, err := GetOrderById(user.ID, orderID)
	if err != nil {
		return nil, err
	}

	// Check if order is paid
	if order.Status == models.OrderStatusWaiting {
		return nil, utils.ErrBadRequest
	}

	var invoiceItems []types.InvoiceItem

	for _, item := range order.OrderItems {
		invoiceItems = append(invoiceItems, types.InvoiceItem{
			Name:     item.Product.Name,
			Weight:   float64(item.Product.Weight),
			Quantity: item.Quantity,
			Price:    utils.FormatToRupiah(item.Product.Price),
			SubTotal: utils.FormatToRupiah(item.SubTotal),
		})
	}

	// Prepare invoice data
	invoice := types.Invoice{
		InvoiceId:     "INV-" + order.ID,
		Date:          order.CreatedAt.Format("02 January 2006"),
		Merchant:      order.Merchant.Name,
		User:          user.Name,
		ContactName:   order.Address.ContactName,
		ContactNumber: order.Address.ContactNumber,
		Address:       ConvertAddressToString(order.Address),
		TotalPrice:    utils.FormatToRupiah(order.TotalPrice),
		PaymentMethod: order.Payment.PaymentMethod,
		Courier:       strings.ToUpper(*order.CourierCompany),
		InvoiceItems:  invoiceItems,
	}

	// Parse template
	tmpl, err := template.New("invoice").Parse(templates.Invoice)
	if err != nil {
		return nil, fmt.Errorf("template parsing error: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, invoice)
	if err != nil {
		return nil, fmt.Errorf("template execution error: %w", err)
	}

	// Create unique temporary files using UUID
	uniqueID := uuid.New().String()
	tempHTML := filepath.Join(os.TempDir(), fmt.Sprintf("invoice_%s_%s.html", orderID, uniqueID))
	tempPDF := filepath.Join(os.TempDir(), fmt.Sprintf("invoice_%s_%s.pdf", orderID, uniqueID))

	// Ensure cleanup of temporary files
	defer func() {
		os.Remove(tempHTML)
		os.Remove(tempPDF)
	}()

	// Write HTML to temp file safely
	err = os.WriteFile(tempHTML, buf.Bytes(), 0600)
	if err != nil {
		return nil, fmt.Errorf("error writing temp HTML file: %w", err)
	}

	// Convert to PDF using wkhtmltopdf
	cmd := exec.Command("wkhtmltopdf", tempHTML, tempPDF)
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("wkhtmltopdf error: %s: %w", string(output), err)
	}

	// Read generated PDF
	pdfBytes, err := os.ReadFile(tempPDF)
	if err != nil {
		return nil, fmt.Errorf("error reading PDF file: %w", err)
	}

	return pdfBytes, nil
}
