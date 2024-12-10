package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/utils"
)

func GetSubCategories(c *gin.Context) {
	page, limit := utils.ExtractPaginationParams(c)

	subCategories, totalItems, err := services.GetSubCategories(page, limit)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error get subcategories")
		return
	}

	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))

	utils.SuccessResponse(c, http.StatusOK, "Success get subcategories", subCategories, page, limit, totalItems, totalPages)
}

func GetSubCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	subCategories, err := services.GetSubCategoryBySlug(slug)

	if err != nil || subCategories == nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Subcategory not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error get subcategory")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success get subcategories", subCategories)
}

func CreateSubCategory(c *gin.Context) {
	var subCategory models.SubCategory

	if !utils.ValidateRequest(c, &subCategory) {
		return
	}

	err := services.CreateSubCategory(&subCategory)

	if err != nil {
		switch err {
		case utils.ErrNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
		case utils.ErrAlreadyExists:
			utils.ErrorResponse(c, http.StatusConflict, "Subcategory name already exists")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error create subcategory")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success create subcategory", subCategory)
}

func UpdateSubCategory(c *gin.Context) {
	var subCategory models.SubCategory
	id := c.Param("id")

	if !utils.ValidateRequest(c, &subCategory) {
		return
	}

	err := services.UpdateSubCategory(id, &subCategory)

	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Subcategory not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error update subcategory")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success update subcategory", subCategory)
}

func DeleteSubCategory(c *gin.Context) {
	id := c.Param("id")

	if err := services.DeleteSubCategory(id); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Subcategory not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error delete subcategory")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success delete subcategory", nil)
}

func UploadSubCategoryIcon(c *gin.Context) {
	filename, err := utils.UploadImageHandler(c, "uploads/subcategory")

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
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error upload subcategory icon")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success upload subcategory icon", filename)
}
