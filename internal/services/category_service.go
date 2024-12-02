package services

import (
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/utils"
	"gorm.io/gorm"
)

func GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := config.DB.Preload("SubCategories").Find(&categories).Error

	return categories, err
}

func GetCategoryBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := config.DB.Preload("SubCategories").First(&category, "slug = ?", slug).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &category, err
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
