package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/utils"
	"gorm.io/gorm"
)

func ClearBrandsCache() error {
	// Get all keys matching the product cache pattern
	pattern := "brands:*"
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

func GetBrands(page, limit int) ([]models.Brand, int64, error) {
	var brands []models.Brand
	var totalCount int64

	cacheKey := fmt.Sprintf("brands:page:%d:limit:%d", page, limit)

	// Check if data is cached
	cachedData, err := config.Rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var cache struct {
			Brands     []models.Brand `json:"brands"`
			TotalCount int64          `json:"totalCount"`
		}

		if err := json.Unmarshal([]byte(cachedData), &cache); err == nil {
			return cache.Brands, cache.TotalCount, nil
		}
	}

	offset := (page - 1) * limit

	if err := config.DB.Model(&models.Brand{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	result := config.DB.Offset(offset).Limit(limit).Find(&brands)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	// Cache the data
	cacheData, err := json.Marshal(struct {
		Brands     []models.Brand `json:"brands"`
		TotalCount int64          `json:"totalCount"`
	}{
		Brands:     brands,
		TotalCount: totalCount,
	})
	if err == nil {
		config.Rdb.Set(context.Background(), cacheKey, string(cacheData), 10*time.Minute).Err()
	}

	return brands, totalCount, nil
}

func GetBrandBySlug(slug string) (*models.Brand, error) {
	var brand models.Brand

	// Check if data is cached
	cacheKey := fmt.Sprintf("brands:slug:%s", slug)
	cachedData, err := config.Rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedData), &brand); err == nil {
			return &brand, nil
		}
	}

	err = config.DB.First(&brand, "slug = ?", slug).Error

	if err == gorm.ErrRecordNotFound {
		return nil, utils.ErrNotFound
	}

	config.Rdb.Set(context.Background(), cacheKey, brand, 10*time.Minute)

	return &brand, nil
}

func CreateBrand(brand *models.Brand) error {
	var existingBrand models.Brand
	err := config.DB.First(&existingBrand, "name = ?", brand.Name).Error

	if err == nil {
		return utils.ErrAlreadyExists
	}

	if err := ClearBrandsCache(); err != nil {
		fmt.Printf("failed to invalidate cache: %v\n", err)
	}

	return config.DB.Create(brand).Error
}

func UpdateBrand(id string, brand *models.Brand) error {
	var existingBrand models.Brand
	err := config.DB.First(&existingBrand, "id = ?", id).Error

	if err == gorm.ErrRecordNotFound {
		return utils.ErrNotFound
	}

	existingBrand.Name = brand.Name
	existingBrand.Icon = brand.Icon

	if err := config.DB.Save(&existingBrand).Error; err != nil {
		return err
	}

	if err := config.DB.First(&existingBrand, "id = ?", id).Error; err != nil {
		return err
	}

	*brand = existingBrand

	if err := ClearBrandsCache(); err != nil {
		fmt.Printf("failed to invalidate cache: %v\n", err)
	}

	return nil
}

func DeleteBrand(id string) error {
	var brand models.Brand

	err := config.DB.First(&brand, "id = ?", id).Error

	if err == gorm.ErrRecordNotFound {
		return utils.ErrNotFound
	}

	// Set slug to null before delete
	config.DB.Model(&models.Brand{}).Where("id = ?", id).Update("slug", nil)

	if err := ClearBrandsCache(); err != nil {
		fmt.Printf("failed to invalidate cache: %v\n", err)
	}

	return config.DB.Delete(&brand).Error
}
