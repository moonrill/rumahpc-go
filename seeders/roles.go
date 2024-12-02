package seeders

import (
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
)

func SeedRoles() {
	roles := []models.Role{
		{Name: "admin"},
		{Name: "customer"},
		{Name: "merchant"},
	}

	for _, role := range roles {
		config.DB.FirstOrCreate(&role, models.Role{Name: role.Name})
	}
}
