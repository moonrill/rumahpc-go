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

func ClearCategoriesCache() error {
	// Get all keys matching the product cache pattern
	pattern := "categories:*"
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

func GetCategories(page, limit int) ([]models.Category, int64, error) {
	var categories []models.Category
	var totalCount int64

	// Check if data is cached
	cacheKey := fmt.Sprintf("categories:page:%d:limit:%d", page, limit)

	cachedData, err := config.Rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var cache struct {
			Categories []models.Category `json:"categories"`
			TotalCount int64             `json:"totalCount"`
		}

		if err := json.Unmarshal([]byte(cachedData), &cache); err == nil {
			return cache.Categories, cache.TotalCount, nil
		}
	}

	offset := (page - 1) * limit

	if err := config.DB.Model(&models.Category{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	result := config.DB.Preload("SubCategories").Offset(offset).Limit(limit).Find(&categories)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	cacheData, err := json.Marshal(struct {
		Categories []models.Category `json:"categories"`
		TotalCount int64             `json:"totalCount"`
	}{
		Categories: categories,
		TotalCount: totalCount,
	})

	if err == nil {
		config.Rdb.Set(context.Background(), cacheKey, string(cacheData), 10*time.Minute).Err()
	}

	return categories, totalCount, nil
}

func GetCategoryBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := config.DB.Preload("SubCategories").First(&category, "slug = ?", slug).Error

	if err == gorm.ErrRecordNotFound {
		return nil, utils.ErrNotFound
	}

	return &category, nil
}

func CreateCategory(category *models.Category) error {
	var existingCategory models.Category
	err := config.DB.First(&existingCategory, "name = ?", category.Name).Error

	if err == nil {
		return utils.ErrAlreadyExists
	}

	err = config.DB.Create(category).Error

	if err != nil {
		return err
	}

	if err := ClearCategoriesCache(); err != nil {
		fmt.Printf("failed to invalidate cache: %v\n", err)
	}

	return nil
}

func UpdateCategory(id string, category *models.Category) error {
	var existingCategory models.Category
	err := config.DB.First(&existingCategory, "id = ?", id).Error

	if err == gorm.ErrRecordNotFound {
		return utils.ErrNotFound
	}

	existingCategory.Name = category.Name
	existingCategory.Icon = category.Icon

	if err := config.DB.Save(&existingCategory).Error; err != nil {
		return err
	}

	if err := config.DB.First(&existingCategory, "id = ?", id).Error; err != nil {
		return err
	}

	*category = existingCategory

	if err := ClearCategoriesCache(); err != nil {
		fmt.Printf("failed to invalidate cache: %v\n", err)
	}

	return nil
}

func DeleteCategory(id string) error {
	var category models.Category
	err := config.DB.First(&category, "id = ?", id).Error

	if err == gorm.ErrRecordNotFound {
		return utils.ErrNotFound
	}

	if err := ClearCategoriesCache(); err != nil {
		fmt.Printf("failed to invalidate cache: %v\n", err)
	}

	return config.DB.Delete(&category).Error
}
