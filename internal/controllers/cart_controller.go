package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/types"
	"github.com/moonrill/rumahpc-api/utils"
)

func AddToCart(c *gin.Context) {
	var requestBody types.AddToCartRequest
	user := c.MustGet("user").(models.User)

	if !utils.ValidateRequest(c, &requestBody) {
		return
	}

	item, err := services.AddToCart(&requestBody, user.ID)

	if err != nil {
		switch err {
		case utils.ErrProductUnavailable:
			utils.ErrorResponse(c, http.StatusUnprocessableEntity, "Product unavailable")
		case utils.ErrNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error add to cart")
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Success add to cart", item)
}

func GetCart(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	cart, err := services.GetCart(user.ID)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error get cart")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success get cart", cart)
}

func UpdateCart(c *gin.Context) {
	id := c.Param("id")
	var requestBody types.UpdateCartRequest
	user := c.MustGet("user").(models.User)

	if !utils.ValidateRequest(c, &requestBody) {
		return
	}

	err := services.UpdateCartItem(id, requestBody.Quantity, user.ID)

	if err != nil {
		switch err {
		case utils.ErrNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Cart item not found")
		case utils.ErrForbidden:
			utils.ErrorResponse(c, http.StatusForbidden, "Forbidden")
		case utils.ErrBadRequest:
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quantity")
		case utils.ErrProductUnavailable:
			utils.ErrorResponse(c, http.StatusUnprocessableEntity, "Product unavailable or stock not enough")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error update cart")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success update cart", nil)
}

func RemoveFromCart(c *gin.Context) {
	var requestBody types.RemoveFromCartRequest
	user := c.MustGet("user").(models.User)

	if !utils.ValidateRequest(c, &requestBody) {
		return
	}

	err := services.RemoveFromCart(requestBody.CartItemsID, user.ID)

	if err != nil {
		switch err {
		case utils.ErrNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Failed to remove from cart", "Some cart item not found")
		case utils.ErrForbidden:
			utils.ErrorResponse(c, http.StatusForbidden, "Failed to remove from cart", "Forbidden")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error remove from cart")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success remove from cart", nil)
}
