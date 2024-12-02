package services

import (
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/utils"
	"gorm.io/gorm"
)

func GetSubCategories() ([]models.SubCategory, error) {
	var subCategories []models.SubCategory
	err := config.DB.Preload("Category").Find(&subCategories).Error

	return subCategories, err
}

func GetSubCategoriesBySlug(slug string) ([]models.SubCategory, error) {
	var subCategories []models.SubCategory
	err := config.DB.Preload("Category").Where("slug = ?", slug).Find(&subCategories).Error

	return subCategories, err
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

	return config.DB.Delete(&subCategory).Error
}
