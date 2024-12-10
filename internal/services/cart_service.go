package services

import (
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/types"
	"github.com/moonrill/rumahpc-api/utils"
	"gorm.io/gorm"
)

func GetGroupedCart(userID string) ([]types.GroupedCartItem, error) {
	var cart models.Cart

	// Preload CartItems, Product, and Merchant (User)
	err := config.DB.
		Where("user_id = ?", userID).
		Preload("CartItems.Product.Merchant").
		First(&cart).Error

	if err != nil {
		return nil, err
	}

	// Map untuk mengelompokkan CartItem berdasarkan MerchantID
	groupedItems := make(map[string]*types.GroupedCartItem)

	// Iterasi setiap CartItem
	for _, item := range cart.CartItems {
		if item.Product != nil && item.Product.Merchant != nil {
			merchantID := item.Product.Merchant.ID
			merchantName := item.Product.Merchant.Name

			// Jika belum ada entry untuk MerchantID ini, buat entry baru
			if _, exists := groupedItems[merchantID]; !exists {
				groupedItems[merchantID] = &types.GroupedCartItem{
					MerchantID:   merchantID,
					MerchantName: merchantName,
					CartItems:    []models.CartItem{},
				}
			}

			// Tambahkan CartItem ke grup merchant
			groupedItems[merchantID].CartItems = append(groupedItems[merchantID].CartItems, item)
		}
	}

	// Ubah map menjadi slice untuk hasil
	var result []types.GroupedCartItem
	for _, group := range groupedItems {
		result = append(result, *group)
	}

	return result, nil
}

func GetCart(userID string) ([]models.CartItem, error) {
	var cart models.Cart

	err := config.DB.Where("user_id = ?", userID).Preload("CartItems").Preload("CartItems.Product").First(&cart).Error

	if err != nil {
		return nil, err
	}

	return cart.CartItems, nil
}

func AddToCart(request *types.AddToCartRequest, userID string) (*models.CartItem, error) {
	// Search product
	var product models.Product

	err := config.DB.Where("id = ?", request.ProductID).First(&product).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	// Check if products is active & stock > 0
	if product.Status != models.ProductStatusActive || product.Stock <= 0 {
		return nil, utils.ErrProductUnavailable
	}

	// Get cart
	var cart models.Cart
	err = config.DB.Where("user_id = ?", userID).First(&cart).Error

	if err != nil {
		return nil, err
	}

	// Check if product is already in cart
	var cartItem models.CartItem
	err = config.DB.Where("cart_id = ? AND product_id = ?", cart.ID, request.ProductID).First(&cartItem).Error

	if err == nil {
		// Add quantity
		cartItem.Quantity += request.Quantity
		err = config.DB.Save(&cartItem).Error
		if err != nil {
			return nil, err
		}
		return &cartItem, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Create new cart item
	cartItem = models.CartItem{
		CartID:    cart.ID,
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
	}

	err = config.DB.Create(&cartItem).Error

	if err != nil {
		return nil, err
	}

	return &cartItem, nil
}

func UpdateCartItem(cartItemID string, quantity int, userID string) error {
	// Validate quantity
	if quantity <= 0 {
		return utils.ErrBadRequest
	}

	// Get cart
	var cart models.Cart
	if err := config.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrNotFound
		}
		return err
	}

	// Get cart item with associated product
	var cartItem models.CartItem
	err := config.DB.Preload("Product").Where("id = ?", cartItemID).First(&cartItem).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrNotFound
		}
		return err
	}

	// Check ownership
	if cartItem.CartID != cart.ID {
		return utils.ErrForbidden
	}

	// Check if product is active
	if cartItem.Product.Status != models.ProductStatusActive {
		return utils.ErrProductUnavailable
	}

	// Check if quantity is more than stock
	if cartItem.Product.Stock < quantity {
		return utils.ErrProductUnavailable
	}

	// Update quantity
	cartItem.Quantity = quantity
	return config.DB.Save(&cartItem).Error
}

func RemoveFromCart(cartItemsID []string, userID string) error {
	// Get cart
	var cart models.Cart
	err := config.DB.Where("user_id = ?", userID).First(&cart).Error

	if err != nil {
		return err
	}

	// Delete cart item
	for _, cartItemID := range cartItemsID {
		var cartItem models.CartItem

		err := config.DB.Where("id = ?", cartItemID).First(&cartItem).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.ErrNotFound
			}
			return err
		}

		if cartItem.CartID != cart.ID {
			return utils.ErrForbidden
		}

		err = config.DB.Delete(&cartItem).Error

		if err != nil {
			return err
		}
	}

	return nil
}
