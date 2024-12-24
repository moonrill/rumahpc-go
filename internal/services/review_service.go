package services

import (
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
)

func GetProductRatingAndReview(productID string) (float64, int64, error) {
	var rating float64
	var reviewCount int64

	err := config.DB.Model(&models.Review{}).
		Where("product_id = ?", productID).
		Select("AVG(rating) as rating, COUNT(*) as review_count").
		Scan(map[string]interface{}{
			"rating":       &rating,
			"review_count": &reviewCount,
		}).Error

	if err != nil {
		return 0, 0, err
	}

	return rating, reviewCount, nil
}
