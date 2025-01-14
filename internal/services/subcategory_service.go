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

func ClearSubCategoriesCache() error {
	// Get all keys matching the product cache pattern
	pattern := "sub_categories:*"
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

func GetSubCategories(page, limit int) ([]models.SubCategory, int64, error) {
	var subCategories []models.SubCategory
	var totalCount int64

	cacheKey := fmt.Sprintf("sub_categories:page:%d:limit:%d", page, limit)

	cachedData, err := config.Rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var cache struct {
			SubCategories []models.SubCategory `json:"sub_categories"`
			TotalCount    int64                `json:"totalCount"`
		}

		if err := json.Unmarshal([]byte(cachedData), &cache); err == nil {
			return cache.SubCategories, cache.TotalCount, nil
		}
	}

	offset := (page - 1) * limit

	if err := config.DB.Model(&models.SubCategory{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	result := config.DB.Preload("Category").Offset(offset).Limit(limit).Find(&subCategories)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	cacheData, err := json.Marshal(struct {
		SubCategories []models.SubCategory `json:"sub_categories"`
		TotalCount    int64                `json:"totalCount"`
	}{
		SubCategories: subCategories,
		TotalCount:    totalCount,
	})

	if err == nil {
		config.Rdb.Set(context.Background(), cacheKey, string(cacheData), 10*time.Minute).Err()
	}

	return subCategories, totalCount, nil
}

func GetSubCategoryBySlug(slug string) (*models.SubCategory, error) {
	var subCategory models.SubCategory
	err := config.DB.Preload("Category").Where("slug = ?", slug).First(&subCategory).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	return &subCategory, nil
}

func CreateSubCategory(subCategory *models.SubCategory) error {
	var category models.Category
	err := config.DB.First(&category, "id = ?", subCategory.CategoryID).Error

	if err != nil {
		return utils.ErrNotFound
	}

	var existingSubCategory models.SubCategory
	err = config.DB.First(&existingSubCategory, "name = ?", subCategory.Name).Error

	if err == nil {
		return utils.ErrAlreadyExists
	}

	err = config.DB.Create(subCategory).Error

	if err != nil {
		return err
	}

	if err := ClearSubCategoriesCache(); err != nil {
		fmt.Printf("failed to invalidate cache: %v\n", err)
	}

	return nil
}

func UpdateSubCategory(id string, subCategory *models.SubCategory) error {
	var existingSubCategory models.SubCategory
	err := config.DB.First(&existingSubCategory, "id = ?", id).Error

	if err == gorm.ErrRecordNotFound {
		return utils.ErrNotFound
	}

	existingSubCategory.Name = subCategory.Name
	existingSubCategory.Icon = subCategory.Icon

	if err := config.DB.Save(&existingSubCategory).Error; err != nil {
		return err
	}

	if err := config.DB.First(&existingSubCategory, "id = ?", id).Error; err != nil {
		return err
	}

	*subCategory = existingSubCategory

	if err := ClearSubCategoriesCache(); err != nil {
		fmt.Printf("failed to invalidate cache: %v\n", err)
	}

	return nil
}

func DeleteSubCategory(id string) error {
	var subCategory models.SubCategory
	err := config.DB.First(&subCategory, "id = ?", id).Error

	if err == gorm.ErrRecordNotFound {
		return utils.ErrNotFound
	}

	// Set slug to null before delete
	config.DB.Model(&models.SubCategory{}).Where("id = ?", id).Update("slug", nil)

	err = config.DB.Delete(&subCategory).Error

	if err != nil {
		return err
	}

	if err := ClearSubCategoriesCache(); err != nil {
		fmt.Printf("failed to invalidate cache: %v\n", err)
	}

	return nil
}
