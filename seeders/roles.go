package seeders

import (
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
)

func SeedRoles() {
	roles := []models.Role{
		{ID: "d9ed6ac7-2a90-4a82-b2a9-c40fda4a5d93", Name: "admin"},
		{ID: "54e26dde-9cc3-45b2-beb5-8b9c02d51eaa", Name: "customer"},
		{ID: "9e1e7628-27c4-49b7-9204-a7b16a009edc", Name: "merchant"},
	}

	for _, role := range roles {
		config.DB.FirstOrCreate(&role, models.Role{Name: role.Name})
	}
}
