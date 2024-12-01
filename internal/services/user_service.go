package services

import (
	"errors"

	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
)

var ErrUserAlreadyExists = errors.New("user already exists")

func CreateUser(user *models.User) error {
	var existingUser models.User
	err := config.DB.First(&existingUser, "email = ?", user.Email).Error

	if err == nil {
		return ErrUserAlreadyExists
	}

	return config.DB.Create(user).Error
}

func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, "email = ?", email).Error

	return &user, err
}
