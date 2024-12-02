package services

import (
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/utils"
	"gorm.io/gorm"
)

func GetBrands(page, limit int) ([]models.Brand, int64, error) {
	var brands []models.Brand
	var totalCount int64

	offset := (page - 1) * limit

	if err := config.DB.Model(&models.Brand{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	result := config.DB.Offset(offset).Limit(limit).Find(&brands)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return brands, totalCount, nil
}

func GetBrandBySlug(slug string) (*models.Brand, error) {
	var brand models.Brand
	err := config.DB.First(&brand, "slug = ?", slug).Error

	if err == gorm.ErrRecordNotFound {
		return nil, utils.ErrNotFound
	}

	return &brand, nil
}

func CreateBrand(brand *models.Brand) error {
	var existingBrand models.Brand
	err := config.DB.First(&existingBrand, "name = ?", brand.Name).Error

	if err == nil {
		return utils.ErrAlreadyExists
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
	existingBrand.Image = brand.Image

	if err := config.DB.Save(&existingBrand).Error; err != nil {
		return err
	}

	if err := config.DB.First(&existingBrand, "id = ?", id).Error; err != nil {
		return err
	}

	*brand = existingBrand

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

	return config.DB.Delete(&brand).Error
}
