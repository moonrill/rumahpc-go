package services

import (
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
)

func GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := config.DB.First(&role, "name = ?", name).Error

	return &role, err
}
