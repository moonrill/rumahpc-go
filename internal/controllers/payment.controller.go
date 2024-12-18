package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/utils"
	"github.com/xendit/xendit-go/v6/invoice"
)

func XenditCallback(c *gin.Context) {
	xtoken := c.GetHeader("x-callback-token")

	if xtoken != os.Getenv("XENDIT_CALLBACK_TOKEN") {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var callback invoice.InvoiceCallback
	if err := c.ShouldBindJSON(&callback); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	err := services.HandleXenditCallback(&callback)

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Callback processed successfully", nil)

}
