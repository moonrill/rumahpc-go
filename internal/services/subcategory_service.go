package services

import (
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/utils"
	"gorm.io/gorm"
)

func GetSubCategories(page, limit int) ([]models.SubCategory, int64, error) {
	var subCategories []models.SubCategory
	var totalCount int64

	offset := (page - 1) * limit

	if err := config.DB.Model(&models.SubCategory{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	result := config.DB.Preload("Category").Offset(offset).Limit(limit).Find(&subCategories)

	if result.Error != nil {
		return nil, 0, result.Error
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

	return config.DB.Create(subCategory).Error
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

	return config.DB.Delete(&subCategory).Error
}
