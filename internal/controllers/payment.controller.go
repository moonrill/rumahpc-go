package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/utils"
)

func XenditCallback(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "success", nil)
}
