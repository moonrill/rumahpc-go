package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/utils"
)

func GetBrands(c *gin.Context) {
	page, limit := utils.ExtractPaginationParams(c)

	brands, totalItems, err := services.GetBrands(page, limit)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error get brands")
		return
	}

	totalPages := int((int64(len(brands)) + int64(limit) - 1) / int64(limit))

	utils.SuccessResponse(c, http.StatusOK, "Success get brands", brands, page, limit, totalItems, totalPages)
}

func GetBrandBySlug(c *gin.Context) {
	slug := c.Param("slug")

	brand, err := services.GetBrandBySlug(slug)

	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Brand not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error get brand")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success get brand", brand)
}

func CreateBrand(c *gin.Context) {
	var brand models.Brand

	if !utils.ValidateRequest(c, &brand) {
		return
	}

	err := services.CreateBrand(&brand)

	if err != nil {
		if err == utils.ErrAlreadyExists {
			utils.ErrorResponse(c, http.StatusConflict, "Brand already exists")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error create brand")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success create brand", brand)
}

func UpdateBrand(c *gin.Context) {
	var brand models.Brand
	id := c.Param("id")

	if !utils.ValidateRequest(c, &brand) {
		return
	}

	err := services.UpdateBrand(id, &brand)

	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Brand not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error update brand")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success update brand", brand)
}

func DeleteBrand(c *gin.Context) {
	id := c.Param("id")

	if err := services.DeleteBrand(id); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Brand not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error delete brand")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success delete brand", nil)
}
