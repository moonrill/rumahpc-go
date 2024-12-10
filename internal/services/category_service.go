package services

import (
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/utils"
	"gorm.io/gorm"
)

func GetCategories(page, limit int) ([]models.Category, int64, error) {
	var categories []models.Category
	var totalCount int64

	offset := (page - 1) * limit

	if err := config.DB.Model(&models.Category{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	result := config.DB.Preload("SubCategories").Offset(offset).Limit(limit).Find(&categories)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return categories, totalCount, nil
}

func GetCategoryBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := config.DB.Preload("SubCategories").First(&category, "slug = ?", slug).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &category, nil
}

func CreateCategory(category *models.Category) error {
	var existingCategory models.Category
	err := config.DB.First(&existingCategory, "name = ?", category.Name).Error

	if err == nil {
		return utils.ErrAlreadyExists
	}

	return config.DB.Create(category).Error
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

	return nil
}

func DeleteCategory(id string) error {
	var category models.Category
	err := config.DB.First(&category, "id = ?", id).Error

	if err == gorm.ErrRecordNotFound {
		return utils.ErrNotFound
	}

	return config.DB.Delete(&category).Error
}
