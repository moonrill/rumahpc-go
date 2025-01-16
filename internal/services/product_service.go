package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/types"
	"github.com/moonrill/rumahpc-api/utils"
	"gorm.io/gorm"
)

func ClearProductsCache() error {
	// Get all keys matching the product cache pattern
	pattern := "products:*"
	ctx := context.Background()

	// Use SCAN to iterate through all keys matching the pattern
	iter := config.Rdb.Scan(ctx, 0, pattern, 0).Iterator()

	// Create a pipeline for batch deletion
	pipe := config.Rdb.Pipeline()

	// Collect all matching keys and delete them
	for iter.Next(ctx) {
		key := iter.Val()
		pipe.Del(ctx, key)
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("error scanning keys: %v", err)
	}

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error clearing cache: %v", err)
	}

	return nil
}

func GetProducts(page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var totalCount int64

	cacheKey := fmt.Sprintf("products:page:%d:limit:%d", page, limit)

	// Check if data is cached
	cachedData, err := config.Rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var cache struct {
			Products   []models.Product `json:"products"`
			TotalCount int64            `json:"totalCount"`
		}

		if err := json.Unmarshal([]byte(cachedData), &cache); err == nil {
			return cache.Products, cache.TotalCount, nil
		}
	}

	offset := (page - 1) * limit

	if err := config.DB.Model(&models.Product{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	result := config.DB.Preload("Brand").Preload("Category").Preload("SubCategory").Offset(offset).Limit(limit).Find(&products)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	for i := range products {
		imageNames, err := GetProductImages(products[i].ID)
		if err != nil {
			return nil, 0, err
		}
		products[i].Images = &imageNames

		rating, reviewCount, err := GetProductRatingAndReview(products[i].ID)
		if err != nil {
			return nil, 0, err
		}
		products[i].Rating = rating
		products[i].ReviewCount = reviewCount
	}

	cacheData, err := json.Marshal(struct {
		Products   []models.Product `json:"products"`
		TotalCount int64            `json:"totalCount"`
	}{
		Products:   products,
		TotalCount: totalCount,
	})
	if err == nil {
		pipe := config.Rdb.Pipeline()

		pipe.Set(context.Background(), cacheKey, cacheData, 10*time.Minute)
		_, _ = pipe.Exec(context.Background())
	}

	return products, totalCount, nil
}

func CreateProduct(product *types.CreateProductRequest, merchantID string) (*models.Product, error) {
	var newProduct models.Product
	newProduct.Name = product.Name
	newProduct.Description = product.Description
	newProduct.Stock = product.Stock
	newProduct.Price = product.Price
	newProduct.Weight = product.Weight
	newProduct.MerchantID = merchantID
	newProduct.BrandID = product.BrandID
	newProduct.CategoryID = product.CategoryID
	newProduct.SubCategoryID = product.SubCategoryID

	if err := config.DB.Create(&newProduct).Error; err != nil {
		return nil, err
	}

	if err := SaveImages(newProduct.ID, product.Images); err != nil {
		return nil, err
	}

	if err := ClearProductsCache(); err != nil {
		fmt.Printf("failed to invalidate cache: %v\n", err)
	}

	return &newProduct, nil
}

func SaveImages(productID string, images []string) error {
	for _, image := range images {
		var productImage models.ProductImages
		productImage.ProductID = productID
		productImage.Image = image
		if err := config.DB.Create(&productImage).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetProductImages(productID string) ([]string, error) {
	var images []string
	err := config.DB.Model(&models.ProductImages{}).Where("product_id = ?", productID).Pluck("image", &images).Error

	if err != nil {
		return nil, err
	}

	return images, nil
}

func GetProductBySlug(slug string) (*models.Product, error) {
	var product models.Product
	if err := config.DB.Preload("Brand").Preload("Category").Preload("SubCategory").Preload("Reviews").Where("slug = ?", slug).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	imageNames, err := GetProductImages(product.ID)
	if err != nil {
		return nil, err
	}
	product.Images = &imageNames

	rating, reviewCount, err := GetProductRatingAndReview(product.ID)
	if err != nil {
		return nil, err
	}
	product.Rating = rating
	product.ReviewCount = reviewCount

	return &product, nil
}

func GetProductsByCategorySlug(slug string, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var totalCount int64

	offset := (page - 1) * limit

	category, err := GetCategoryBySlug(slug)

	if err != nil {
		return nil, 0, utils.ErrNotFound
	}

	if err := config.DB.Model(&models.Product{}).Where("category_id = ?", category.ID).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	result := config.DB.Preload("Brand").Preload("Category").Preload("SubCategory").Offset(offset).Limit(limit).Where("category_id = ?", category.ID).Find(&products)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	for i := range products {
		imageNames, err := GetProductImages(products[i].ID)
		if err != nil {
			return nil, 0, err
		}
		products[i].Images = &imageNames

		rating, reviewCount, err := GetProductRatingAndReview(products[i].ID)
		if err != nil {
			return nil, 0, err
		}
		products[i].Rating = rating
		products[i].ReviewCount = reviewCount
	}

	return products, totalCount, nil
}

func GetProductsBySubCategorySlug(slug string, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var totalCount int64

	offset := (page - 1) * limit

	subCategory, err := GetSubCategoryBySlug(slug)

	if err != nil {
		return nil, 0, utils.ErrNotFound
	}

	if err := config.DB.Model(&models.Product{}).Where("sub_category_id = ?", subCategory.ID).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	result := config.DB.Preload("Brand").Preload("Category").Preload("SubCategory").Offset(offset).Limit(limit).Where("sub_category_id = ?", subCategory.ID).Find(&products)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	for i := range products {
		imageNames, err := GetProductImages(products[i].ID)
		if err != nil {
			return nil, 0, err
		}
		products[i].Images = &imageNames

		rating, reviewCount, err := GetProductRatingAndReview(products[i].ID)
		if err != nil {
			return nil, 0, err
		}
		products[i].Rating = rating
		products[i].ReviewCount = reviewCount
	}

	return products, totalCount, nil
}

func UpdateProduct(id string, product *types.UpdateProductRequest, merchantID string) (*models.Product, error) {
	var updatedProduct models.Product

	if err := config.DB.Where("id = ?", id).First(&updatedProduct).Error; err != nil {
		return nil, utils.ErrNotFound
	}

	if updatedProduct.MerchantID != merchantID {
		return nil, utils.ErrForbidden
	}

	updatedProduct.Name = product.Name
	updatedProduct.Description = product.Description
	updatedProduct.Stock = product.Stock
	updatedProduct.Price = product.Price
	updatedProduct.Weight = product.Weight
	updatedProduct.BrandID = product.BrandID
	updatedProduct.CategoryID = product.CategoryID
	updatedProduct.SubCategoryID = product.SubCategoryID

	if err := config.DB.Save(&updatedProduct).Error; err != nil {
		return nil, err
	}

	if err := ClearProductsCache(); err != nil {
		fmt.Printf("failed to invalidate cache: %v\n", err)
	}

	return &updatedProduct, nil
}

func ToggleProductStatus(id string, merchantID string) error {
	var product models.Product

	if err := config.DB.Where("id = ?", id).First(&product).Error; err != nil {
		return utils.ErrNotFound
	}

	if product.MerchantID != merchantID {
		return utils.ErrForbidden
	}

	if product.Status == models.ProductStatusActive {
		product.Status = models.ProductStatusInactive
	} else {
		product.Status = models.ProductStatusActive
	}

	if err := config.DB.Save(&product).Error; err != nil {
		return err
	}

	if err := ClearProductsCache(); err != nil {
		fmt.Printf("failed to invalidate cache: %v\n", err)
	}

	return nil
}

func GetProductRecommendations(slug string, limit int) ([]models.Product, error) {
	var product models.Product
	var products []models.Product

	// First find the source product
	if err := config.DB.Where("slug = ?", slug).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	cacheKey := fmt.Sprintf("product:recommendations:%s", slug)
	cachedData, err := config.Rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedData), &products); err == nil {
			return products, nil
		}
	}

	// Split the product name and description into keywords
	nameKeywords := strings.Fields(strings.ToLower(product.Name))
	descKeywords := strings.Fields(strings.ToLower(product.Description))

	// Create a query builder for primary recommendations
	query := config.DB.Model(&models.Product{}).
		Where("status = ?", models.ProductStatusActive).
		Where("id != ?", product.ID)

	// Build dynamic OR conditions for name matching
	var nameConditions []string
	var args []interface{}

	// Add exact category matching if your product model has categories
	if product.CategoryID != "" {
		query = query.Where("category_id = ?", product.CategoryID)
	}
	if product.SubCategoryID != nil {
		query = query.Where("sub_category_id = ?", product.SubCategoryID)
	}

	// Add conditions for name keywords
	for _, keyword := range nameKeywords {
		if len(keyword) > 3 { // Only use keywords longer than 3 characters
			nameConditions = append(nameConditions, "name ILIKE ?")
			args = append(args, "%"+keyword+"%")
		}
	}

	// Add conditions for description keywords
	for _, keyword := range descKeywords {
		if len(keyword) > 3 { // Only use keywords longer than 3 characters
			nameConditions = append(nameConditions, "description ILIKE ?")
			args = append(args, "%"+keyword+"%")
		}
	}

	// Combine all conditions
	if len(nameConditions) > 0 {
		query = query.Where(strings.Join(nameConditions, " OR "), args...)
	}

	// Add relevance scoring
	query = query.Select("products.*, "+
		"CASE "+
		"WHEN category_id = ? AND sub_category_id = ? THEN 4 "+
		"WHEN category_id = ? THEN 3 "+
		"WHEN sub_category_id = ? THEN 2 "+
		"ELSE 1 END as relevance_score",
		product.CategoryID, product.SubCategoryID,
		product.CategoryID,
		product.SubCategoryID)

	// Get primary recommendations
	err = query.Order("relevance_score DESC").
		Limit(limit).
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	// If we don't have enough recommendations, get additional similar products
	if len(products) < limit {
		remainingLimit := limit - len(products)
		var additionalProducts []models.Product

		// Create a new query for additional products
		fallbackQuery := config.DB.Model(&models.Product{}).
			Where("status = ?", models.ProductStatusActive).
			Where("id != ?", product.ID).
			// Exclude already recommended products
			Where("id NOT IN (?)", getProductIDs(products))

		// Try to find products in the same category first
		if product.CategoryID != "" {
			fallbackQuery = fallbackQuery.Where("category_id = ?", product.CategoryID)
		}

		err = fallbackQuery.
			Select("products.*, 1 as relevance_score").
			Order("RANDOM()"). // Add some randomness to recommendations
			Limit(remainingLimit).
			Find(&additionalProducts).Error

		if err != nil {
			return nil, err
		}

		// If we still don't have enough, try without category restriction
		if len(additionalProducts) < remainingLimit {
			var finalProducts []models.Product
			finalLimit := remainingLimit - len(additionalProducts)

			err = config.DB.Model(&models.Product{}).
				Where("status = ?", models.ProductStatusActive).
				Where("id != ?", product.ID).
				Where("id NOT IN (?)", getProductIDs(append(products, additionalProducts...))).
				Select("products.*, 0 as relevance_score").
				Order("RANDOM()").
				Limit(finalLimit).
				Find(&finalProducts).Error

			if err != nil {
				return nil, err
			}

			additionalProducts = append(additionalProducts, finalProducts...)
		}

		// Combine primary and additional recommendations
		products = append(products, additionalProducts...)
	}

	// Cache the results
	cacheData, err := json.Marshal(products)
	if err == nil {
		config.Rdb.Set(context.Background(), cacheKey, cacheData, 10*time.Minute).Err()
	}

	return products, nil
}

// Helper function to extract product IDs
func getProductIDs(products []models.Product) []string {
	ids := make([]string, len(products))
	for i, p := range products {
		ids[i] = p.ID
	}
	return ids
}
