package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/utils"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var request models.User

	if !utils.ValidateRequest(c, &request) {
		return
	}

	hashedPassword, salt, err := services.HashPassword(request.Password)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error hashing password", err.Error())
		return
	}

	user := models.User{
		Name:        request.Name,
		Email:       request.Email,
		Password:    hashedPassword,
		PhoneNumber: request.PhoneNumber,
		Salt:        salt,
	}

	err = services.CreateUser(&user)

	if err != nil {
		if err == services.ErrUserAlreadyExists {
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
	c.SetCookie("access_token", token, 3600*24*7, "", "", false, true)

	utils.SuccessResponse(c, http.StatusOK, "Success login", gin.H{"access_token": token})
}
