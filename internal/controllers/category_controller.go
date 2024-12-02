package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/utils"
)

func GetCategories(c *gin.Context) {
	categories, err := services.GetCategories()

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error get categories")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success get categories", categories)
}

func GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	category, err := services.GetCategoryBySlug(slug)
	if err != nil || category == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success get category", category)
}

func CreateCategory(c *gin.Context) {
	var category models.Category

	if !utils.ValidateRequest(c, &category) {
		return
	}

	err := services.CreateCategory(&category)

	if err != nil {
		if err == utils.ErrAlreadyExists {
			utils.ErrorResponse(c, http.StatusConflict, "Category name already exists")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error create category")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success create category", category)
}

func UpdateCategory(c *gin.Context) {
	var category models.Category
	id := c.Param("id")

	if !utils.ValidateRequest(c, &category) {
		return
	}

	err := services.UpdateCategory(id, &category)

	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error update category")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success update category", category)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	if err := services.DeleteCategory(id); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error delete category")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success delete category", nil)
}
