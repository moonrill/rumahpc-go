package utils

import "github.com/gin-gonic/gin"

func SuccessResponse(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, gin.H{"data": data, "status_code": status, "message": message})
}

func ErrorResponse(c *gin.Context, status int, message string, err ...any) {
	response := gin.H{"status_code": status, "message": message}

	if len(err) > 0 {
		response["error"] = err[0]
	}

	c.JSON(status, response)
}
