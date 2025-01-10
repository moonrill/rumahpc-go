package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/utils"
)

func GeneratePC(c *gin.Context) {
	computer, err := services.GeneratePC()

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error generate pc", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success generate pc", computer)
}
