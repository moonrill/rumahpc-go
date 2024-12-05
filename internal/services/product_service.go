package services

import (
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/types"
	"github.com/moonrill/rumahpc-api/utils"
	"gorm.io/gorm"
)

func GetProducts(page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var totalCount int64

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
	if err := config.DB.Preload("Brand").Preload("Category").Preload("SubCategory").Where("slug = ?", slug).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
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

	return nil
}
