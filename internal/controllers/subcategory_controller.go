package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/utils"
)

func GetSubCategoriesByCategoryID(c *gin.Context) {
	id := c.Param("id")
	subCategories, err := services.GetSubCategoriesByCategoryID(id)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error get subcategories")
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
		case services.ErrCategoryNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
		case services.ErrSubCategoryAlreadyExists:
			utils.ErrorResponse(c, http.StatusConflict, "Subcategory name already exists")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error create subcategory")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success create subcategory", subCategory)
}
