package services

import (
	"errors"

	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
)

var ErrSubCategoryAlreadyExists = errors.New("subcategory name already exists")

func GetSubCategoriesByCategoryID(id string) ([]models.SubCategory, error) {
	var subCategories []models.SubCategory
	err := config.DB.Preload("Category").Where("category_id = ?", id).Find(&subCategories).Error

	return subCategories, err
}

func CreateSubCategory(subCategory *models.SubCategory) error {
	var category models.Category
	err := config.DB.First(&category, "id = ?", subCategory.CategoryID).Error

	if err != nil {
		return ErrCategoryNotFound
	}

	var existingSubCategory models.SubCategory
	err = config.DB.First(&existingSubCategory, "name = ?", subCategory.Name).Error

	if err == nil {
		return ErrSubCategoryAlreadyExists
	}

	return config.DB.Create(subCategory).Error
}
