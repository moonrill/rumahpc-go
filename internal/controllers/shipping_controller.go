package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/types"
	"github.com/moonrill/rumahpc-api/utils"
)

func GetBuyNowCouriersRates(c *gin.Context) {
	var requestBody types.BuyNowCouriersRatesRequest

	if !utils.ValidateRequest(c, &requestBody) {
		return
	}

	data, err := services.GetBuyNowCouriersRates(&requestBody)

	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Address not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error get couriers rates")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success get couriers rates", data)
}

func GetCartCouriersRates(c *gin.Context) {
	var requestBody types.CartCouriersRatesRequest

	if !utils.ValidateRequest(c, &requestBody) {
		return
	}

	data, err := services.GetCartCouriersRates(&requestBody)

	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Address not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error get couriers rates")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success get couriers rates", data)
}

func BiteshipCallback(c *gin.Context) {
	signature := c.GetHeader("x-biteship-signature")

	if signature != os.Getenv("BITESHIP_CALLBACK_SIGNATURE") {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var callback types.BiteshipStatusCallback
	if err := c.ShouldBindJSON(&callback); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	err := services.HandleBiteshipCallback(&callback)

	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error get order")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Webhook processed successfully", nil)
}
