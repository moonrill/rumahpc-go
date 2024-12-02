package utils

import "github.com/gin-gonic/gin"

func SuccessResponse(c *gin.Context, status int, message string, data interface{}, another ...any) {
	response := gin.H{"data": data, "status_code": status, "message": message}

	if len(another) >= 4 {
		response["page"] = another[0]
		response["limit"] = another[1]
		response["total_items"] = another[2]
		response["total_pages"] = another[3]
	}

	c.JSON(status, response)
}

func ErrorResponse(c *gin.Context, status int, message string, err ...any) {
	response := gin.H{"status_code": status, "message": message}

	if len(err) > 0 {
		response["error"] = err[0]
	}

	c.JSON(status, response)
}
