package services

import (
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/types"
	"github.com/moonrill/rumahpc-api/utils"
	"gorm.io/gorm"
)

func GetProductRatingAndReview(productID string) (float64, int64, error) {
	var result struct {
		Rating      float64
		ReviewCount int64
	}

	err := config.DB.Model(&models.Review{}).
		Select("AVG(rating) as rating, COUNT(*) as review_count").
		Where("product_id = ?", productID).
		Scan(&result).Error

	if err != nil {
		return 0, 0, err
	}

	return result.Rating, result.ReviewCount, nil
}

func CreateReview(request *types.ReviewRequest, userID string) (*models.Review, error) {
	var order models.Order
	var orderItem models.OrderItem
	var existingReview models.Review

	if err := config.DB.First(&order, "id = ? AND user_id = ?", request.OrderID, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	if order.Status != models.OrderStatusCompleted {
		return nil, utils.ErrBadRequest
	}

	if err := config.DB.Preload("Product.Merchant").First(&orderItem, "order_id = ? AND product_id = ?", request.OrderID, request.ProductID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	if err := config.DB.First(&existingReview, "order_item_id = ? AND product_id = ?", orderItem.ID, request.ProductID).Error; err == nil {
		return nil, utils.ErrReviewAlreadyExists
	}

	newReview := &models.Review{
		UserID:      userID,
		MerchantID:  orderItem.Product.MerchantID,
		ProductID:   request.ProductID,
		OrderItemID: orderItem.ID,
		Rating:      request.Rating,
		Comment:     &request.Comment,
	}

	if err := config.DB.Create(newReview).Error; err != nil {
		return nil, err
	}

	return newReview, nil
}

func GetUserReviews(userId, sort string, page, limit int) ([]*models.Review, int64, error) {
	var reviews []*models.Review
	var totalCount int64

	offset := (page - 1) * limit

	if err := config.DB.Model(&models.Review{}).Where("user_id = ?", userId).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := config.DB.Preload("OrderItem.Product").Where("user_id = ?", userId).Offset(offset).Limit(limit).Order("created_at " + sort).Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	for i := range reviews {
		productImageNames, err := GetProductImages(reviews[i].OrderItem.Product.ID)
		if err != nil {
			return nil, 0, err
		}

		reviews[i].OrderItem.Product.Images = &productImageNames
	}

	return reviews, totalCount, nil
}
