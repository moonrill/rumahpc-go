package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateRequest(c *gin.Context, body interface{}) bool {
	if err := c.ShouldBindJSON(body); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return false
	}

	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				errors = append(errors, err.Field()+" is required")
			case "email":
				errors = append(errors, err.Field()+" must be a valid email")
			case "min":
				errors = append(errors, err.Field()+" must be at least "+err.Param()+" characters")
			case "max":
				errors = append(errors, err.Field()+" must be at most "+err.Param()+" characters")
			case "gte":
				errors = append(errors, err.Field()+" must be greater than or equal to "+err.Param())
			case "lte":
				errors = append(errors, err.Field()+" must be less than or equal to "+err.Param())
			case "containsany":
				errors = append(errors, err.Field()+" must contain at least one special character")
			default:
				errors = append(errors, "Validation error for "+err.Field())
			}
		}

		ErrorResponse(c, http.StatusBadRequest, "Validation error", errors)
		return false
	}

	return true
}
