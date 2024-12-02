package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/internal/services"
	"github.com/moonrill/rumahpc-api/utils"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Email       string `json:"email" validate:"required,email,max=255"`
	Password    string `json:"password" validate:"required,min=8,max=255"`
	PhoneNumber string `json:"phone_number" validate:"required,max=13"`
	Role        string `json:"role" validate:"required,max=255"`
}

func SignUp(c *gin.Context) {
	var request SignUpRequest

	if !utils.ValidateRequest(c, &request) {
		return
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
		RoleID:      role.ID,
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
	c.SetCookie("Authorization", token, 3600*24*7, "", "", false, true)

	utils.SuccessResponse(c, http.StatusOK, "Success login", gin.H{"access_token": token})
}

func GetProfile(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	utils.SuccessResponse(c, http.StatusOK, "Success get profile", user)
}
