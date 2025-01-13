package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/types"
	"github.com/moonrill/rumahpc-api/utils"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

	otp := utils.GenerateOTP(6)
	otpKey := "otp_" + user.ID
	lastRequestKey := "last_otp_request_" + user.ID
	err = config.Rdb.Set(context.Background(), otpKey, otp, 5*time.Minute).Err()

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error set otp", err.Error())
		return
	}

	err = config.Rdb.Set(context.Background(), lastRequestKey, time.Now(), time.Minute).Err()

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error set last request", err.Error())
		return
	}

	err = utils.SendOTPEmail(user.Email, otp)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error send otp", err.Error())
		return
	}

	response := map[string]interface{}{
		"user":           user,
		"otp_expired_at": time.Now().Add(5 * time.Minute).Format("2006-01-02T15:04:05Z07:00"),
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success create user", response)
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

	if !user.IsActive {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User is not active")
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

func VerifyOTP(c *gin.Context) {
	var request types.OTPRequest
	if !utils.ValidateRequest(c, &request) {
		return
	}

	otpKey := "otp_" + request.UserID
	storedOtp, err := config.Rdb.Get(context.Background(), otpKey).Result()

	if err == redis.Nil {
		utils.ErrorResponse(c, http.StatusNotFound, "OTP Expired or Invalid")
		return
	} else if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error get otp", err.Error())
		return
	}

	if storedOtp != request.OTP {
		utils.ErrorResponse(c, http.StatusNotFound, "Invalid OTP")
		return
	}

	err = services.ActivateUser(request.UserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error activate user", err.Error())
		return
	}

	config.Rdb.Del(context.Background(), otpKey)

	utils.SuccessResponse(c, http.StatusOK, "Success activate user", nil)
}

func ResendOTP(c *gin.Context) {
	var request types.ResendOTPRequest
	var user models.User

	if !utils.ValidateRequest(c, &request) {
		return
	}

	err := config.DB.First(&user, "id = ?", request.UserID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error get user", err.Error())
		}
		return
	}

	if user.IsActive {
		utils.ErrorResponse(c, http.StatusBadRequest, "User is already active")
		return
	}

	// Check if the user has already requested an OTP in the last minute
	lastRequestKey := "last_otp_request_" + user.ID
	lastRequest, err := config.Rdb.Get(context.Background(), lastRequestKey).Result()

	if lastRequest != "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "You have already requested an OTP in the last minute")
		return
	}

	if err != redis.Nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error get last request", err.Error())
		return
	}

	err = config.Rdb.Set(context.Background(), lastRequestKey, time.Now(), time.Minute).Err()

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error set last request", err.Error())
		return
	}

	otp := utils.GenerateOTP(6)
	otpKey := "otp_" + user.ID
	err = config.Rdb.Set(context.Background(), otpKey, otp, 5*time.Minute).Err()

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error set otp", err.Error())
		return
	}

	err = utils.SendOTPEmail(user.Email, otp)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error send otp", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success resend otp", nil)
}
