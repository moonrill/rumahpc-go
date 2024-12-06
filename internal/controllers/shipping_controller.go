package controllers

import (
	"net/http"

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
