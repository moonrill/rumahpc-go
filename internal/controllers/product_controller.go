package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/types"
	"github.com/moonrill/rumahpc-api/utils"
)

func GetProducts(c *gin.Context) {
	page, limit := utils.ExtractPaginationParams(c)

	products, totalItems, err := services.GetProducts(page, limit)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error get products")
		return
	}

	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))

	utils.SuccessResponse(c, http.StatusOK, "Success get products", products, page, limit, totalItems, totalPages)
}

func GetProduct(c *gin.Context) {
	slug := c.Param("slug")
	product, err := services.GetProductBySlug(slug)

	if err != nil || product == nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error get product")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success get product", product)
}

func CreateProduct(c *gin.Context) {
	merchant := c.MustGet("user").(models.User)
	var requestBody types.CreateProductRequest

	if !utils.ValidateRequest(c, &requestBody) {
		return
	}

	product, err := services.CreateProduct(&requestBody, merchant.ID)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error create product")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success create product", product)
}

func GetProductsByCategorySlug(c *gin.Context) {
	cSlug := c.Param("slug")

	page, limit := utils.ExtractPaginationParams(c)

	products, totalItems, err := services.GetProductsByCategorySlug(cSlug, page, limit)

	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error get products")
		}
		return
	}

	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))

	utils.SuccessResponse(c, http.StatusOK, "Success get products", products, page, limit, totalItems, totalPages)
}

func GetProductsBySubCategorySlug(c *gin.Context) {
	scSlug := c.Param("slug")

	page, limit := utils.ExtractPaginationParams(c)

	products, totalItems, err := services.GetProductsBySubCategorySlug(scSlug, page, limit)

	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Sub Category not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error get products")
		}
		return
	}

	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))

	utils.SuccessResponse(c, http.StatusOK, "Success get products", products, page, limit, totalItems, totalPages)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	merchant := c.MustGet("user").(models.User)

	var requestBody types.UpdateProductRequest

	if !utils.ValidateRequest(c, &requestBody) {
		return
	}

	updatedProduct, err := services.UpdateProduct(id, &requestBody, merchant.ID)

	if err != nil {
		switch err {
		case utils.ErrNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
		case utils.ErrForbidden:
			utils.ErrorResponse(c, http.StatusForbidden, "Forbidden")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error update product")
		}

		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success update product", updatedProduct)
}

func ToggleProductStatus(c *gin.Context) {
	id := c.Param("id")
	merchant := c.MustGet("user").(models.User)

	err := services.ToggleProductStatus(id, merchant.ID)

	if err != nil {
		switch err {
		case utils.ErrNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
		case utils.ErrForbidden:
			utils.ErrorResponse(c, http.StatusForbidden, "Forbidden")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error toggle product status")
		}

		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success toggle product status", nil)
}

func UploadProductImage(c *gin.Context) {
	filename, err := utils.UploadImageHandler(c, "uploads/product")

	if err != nil {
		switch err {
		case utils.ErrUploadImage:
			utils.ErrorResponse(c, http.StatusBadRequest, "Failed to upload image")
		case utils.ErrUploadImageExt:
			utils.ErrorResponse(c, http.StatusBadRequest, "Image extension not allowed")
		case utils.ErrUploadImageSize:
			utils.ErrorResponse(c, http.StatusBadRequest, "Image size too large")
		case utils.ErrSaveImage:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save image")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error upload product image")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success upload product image", filename)
}

func UploadMultipleProductImages(c *gin.Context) {
	files, err := utils.UploadMultipleImageHandler(c, "uploads/product")

	if err != nil {
		switch err {
		case utils.ErrUploadImage:
			utils.ErrorResponse(c, http.StatusBadRequest, "Failed to upload image")
		case utils.ErrUploadImageExt:
			utils.ErrorResponse(c, http.StatusBadRequest, "File extension not allowed")
		case utils.ErrUploadImageSize:
			utils.ErrorResponse(c, http.StatusBadRequest, "Image size too large")
		case utils.ErrSaveImage:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save image")
		case utils.ErrEmptyUpload:
			utils.ErrorResponse(c, http.StatusBadRequest, "No files provided")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error upload product image")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success upload product image", files)
}
