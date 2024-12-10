package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ExtractPaginationParams(c *gin.Context) (int, int) {
	defaultPage := 1
	defaultLimit := 10

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))

	if err != nil || page < 1 {
		page = defaultPage
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err != nil || limit < 1 {
		limit = defaultLimit
	}

	return page, limit
}
