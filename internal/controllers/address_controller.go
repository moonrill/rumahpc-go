package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/utils"
	"gorm.io/gorm"
)

func CreateAddress(c *gin.Context) {
	var address models.Address

	if !utils.ValidateRequest(c, &address) {
		return
	}

	// If user role is merchant the address must be 1
	user := c.MustGet("user").(models.User)

	if user.Role.Name == "merchant" {
		var existAdress models.Address

		err := config.DB.First(&existAdress, "user_id = ?", user.ID).Error

		if err == nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Merchant only can have 1 address")
			return
		}
	}

	if address.Default {
		config.DB.Model(&models.Address{}).Where("user_id = ?", user.ID).Update("default", false)
	}

	address.UserID = user.ID

	err := config.DB.Create(&address).Error

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error create address")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success create address", address)
}

func GetAddress(c *gin.Context) {
	var addresses []models.Address

	user := c.MustGet("user").(models.User)

	err := config.DB.Find(&addresses, "user_id = ?", user.ID).Error

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error get address")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success get address", addresses)
}

func GetAddressById(c *gin.Context) {
	id := c.Param("id")

	var address models.Address

	user := c.MustGet("user").(models.User)

	err := config.DB.First(&address, "id = ? AND user_id = ?", id, user.ID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Address not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error get address")
		}

		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success get address", address)
}

func UpdateAddress(c *gin.Context) {
	id := c.Param("id")

	user := c.MustGet("user").(models.User)

	var requestBody models.Address

	if !utils.ValidateRequest(c, &requestBody) {
		return
	}

	var existingAddress models.Address

	err := config.DB.First(&existingAddress, "id = ? AND user_id = ?", id, user.ID).Error

	if err == gorm.ErrRecordNotFound {
		utils.ErrorResponse(c, http.StatusNotFound, "Address not found")
		return
	}

	existingAddress.ContactName = requestBody.ContactName
	existingAddress.ContactNumber = requestBody.ContactNumber
	existingAddress.Province = requestBody.Province
	existingAddress.City = requestBody.City
	existingAddress.District = requestBody.District
	existingAddress.Village = requestBody.Village
	existingAddress.ZipCode = requestBody.ZipCode
	existingAddress.Address = requestBody.Address
	existingAddress.Note = requestBody.Note
	existingAddress.Default = requestBody.Default

	if existingAddress.Default {
		config.DB.Model(&models.Address{}).Where("user_id = ?", user.ID).Update("default", false)
	}

	err = config.DB.Save(&existingAddress).Error

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error update address")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success update address", existingAddress)
}

func DeleteAddress(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Address id is required")
		return
	}

	user := c.MustGet("user").(models.User)

	var address models.Address

	err := config.DB.First(&address, "id = ? AND user_id = ?", id, user.ID).Error

	if err == gorm.ErrRecordNotFound {
		utils.ErrorResponse(c, http.StatusNotFound, "Address not found")
		return
	}

	err = config.DB.Delete(&address).Error

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error delete address")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success delete address", nil)
}
