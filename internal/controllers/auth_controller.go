package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/types"
	"github.com/moonrill/rumahpc-api/utils"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var request types.SignUpRequest

	if !utils.ValidateRequest(c, &request) {
		return
	}

	if request.Role == "merchant" {
		var errors []string
		if request.PaymentChannel == "" {
			errors = append(errors, "Payment channel is required")
		}
		if request.AccountNumber == "" {
			errors = append(errors, "Account number is required")
		}
		if request.AccountName == "" {
			errors = append(errors, "Account name is required")
		}

		if len(errors) > 0 {
			utils.ErrorResponse(c, http.StatusBadRequest, "Validation error", errors)
			return
		}
	}

	role, err := services.GetRoleByName(request.Role)

	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Role not found")
		return
	}

	if role.Name == "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "Admin role not allowed")
		return
	}

	// Check if email already exists
	if services.CheckEmailExists(request.Email) {
		utils.ErrorResponse(c, http.StatusConflict, "Email already used")
		return
	}

	hashedPassword, salt, err := services.HashPassword(request.Password)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error hashing password", err.Error())
		return
	}

	user := models.User{
		ID:          uuid.New().String(),
		Name:        request.Name,
		Email:       request.Email,
		Password:    hashedPassword,
		PhoneNumber: request.PhoneNumber,
		Salt:        salt,
		RoleID:      role.ID,
	}

	err = services.CreateUser(&user)

	// Create Cart for Customer
	if role.Name == "customer" {
		cart := models.Cart{
			UserID: user.ID,
		}
		err = config.DB.Create(&cart).Error

		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error create cart", err.Error())
			return
		}
	}

	if err != nil {
		if err == utils.ErrAlreadyExists {
			utils.ErrorResponse(c, http.StatusConflict, "User already exists")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error create user", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success create user", user)
}

func SignIn(c *gin.Context) {
	var request struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if !utils.ValidateRequest(c, &request) {
		return
	}

	user, err := services.FindUserByEmail(request.Email)

	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password+user.Salt))

	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := services.GenerateToken(user)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error generating token", err.Error())
		return
	}

	// Set the token as a cookie in the response
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600*24*7, "", "", false, true)

	utils.SuccessResponse(c, http.StatusOK, "Success login", gin.H{"access_token": token})
}

func GetProfile(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	utils.SuccessResponse(c, http.StatusOK, "Success get profile", user)
}
