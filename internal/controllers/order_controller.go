package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/types"
	"github.com/moonrill/rumahpc-api/utils"
)

func BuyNowOrder(c *gin.Context) {
	var request types.BuyNowRequest
	user := c.MustGet("user").(models.User)

	if !utils.ValidateRequest(c, &request) {
		return
	}

	order, err := services.CreateBuyNowOrder(&request, user.ID)

	if err != nil {
		switch err {
		case services.ErrOrderProductNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
		case services.ErrOrderAddressNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Address not found")
		case utils.ErrProductUnavailable:
			utils.ErrorResponse(c, http.StatusUnprocessableEntity, "Product unavailable")
		case services.ErrOrderShipping:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error shipping order")
		case services.ErrOrderCreate:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error creating order")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}

		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success creating order", order)
}

func CheckoutCart(c *gin.Context) {
	var request types.CheckoutCartRequest
	user := c.MustGet("user").(models.User)

	if !utils.ValidateRequest(c, &request) {
		return
	}

	invoice, err := services.CreateCartCheckoutOrder(&request, &user)

	if err != nil {
		switch err {
		case services.ErrOrderProductNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
		case services.ErrOrderAddressNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Address not found")
		case services.ErrOrderCartItemNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "One or more cart items not found")
		case services.ErrOrderShipping:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error shipping order")
		case services.ErrOrderCreate:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error creating order")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}

		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success creating order", invoice)
}

func GetOrders(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	page, limit := utils.ExtractPaginationParams(c)
	status := c.Query("status")
	shippingStatus := c.Query("shipping_status")

	orders, totalItems, err := services.GetOrders(user.ID, page, limit, status, shippingStatus)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error get orders")
		return
	}

	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))

	utils.SuccessResponse(c, http.StatusOK, "Success get orders", orders, page, limit, totalItems, totalPages)
}

func GetOrderById(c *gin.Context) {
	orderID := c.Param("id")
	user := c.MustGet("user").(models.User)

	order, err := services.GetOrderById(user.ID, orderID)

	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error get order")
		}

		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success get order", order)
}

func CompleteOrder(c *gin.Context) {
	orderID := c.Param("id")
	user := c.MustGet("user").(models.User)

	err := services.CompleteOrder(orderID, user.ID)

	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
		} else if err == utils.ErrBadRequest {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order status")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error complete order")
		}

		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success complete order", nil)
}
