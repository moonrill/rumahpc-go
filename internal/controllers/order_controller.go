package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/types"
	"github.com/moonrill/rumahpc-api/utils"
)

func CreateBuyNowOrder(c *gin.Context) {
	var request types.BuyNowRequest
	user := c.MustGet("user").(models.User)

	if !utils.ValidateRequest(c, &request) {
		return
	}

	order, err := services.CreateBuyNowOrder(&request, user.ID)

	if err != nil {
		switch err {
		case services.ErrOrderProductNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
		case services.ErrOrderAddressNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Address not found")
		case services.ErrOrderShipping:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error shipping order")
		case services.ErrOrderCreate:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error creating order")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}

		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success creating order", order)
}

func CreateCartOrder(c *gin.Context) {

}
