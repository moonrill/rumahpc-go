package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/types"
	"github.com/moonrill/rumahpc-api/utils"
)

func CreateReview(c *gin.Context) {
	var request types.ReviewRequest
	user := c.MustGet("user").(models.User)

	if !utils.ValidateRequest(c, &request) {
		return
	}

	result, err := services.CreateReview(&request, user.ID)

	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Order or Order Item not found")
		} else if err == utils.ErrReviewAlreadyExists {
			utils.ErrorResponse(c, http.StatusUnprocessableEntity, "Review already exists")
		} else if err == utils.ErrBadRequest {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order status")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error create review", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success create review", result)
}

func GetUserReviews(c *gin.Context) {
	sort := c.Query("sort")
	user := c.MustGet("user").(models.User)
	page, limit := utils.ExtractPaginationParams(c)

	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	reviews, totalItems, err := services.GetUserReviews(user.ID, sort, page, limit)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error get user reviews", err.Error())
		return
	}

	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))

	utils.SuccessResponse(c, http.StatusOK, "Success get user reviews", reviews, page, limit, totalItems, totalPages)
}
