package seeders

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/internal/services"
)

func SeedUsers() {
	users := []models.User{
		{
			Name:        "RumahPC Admin",
			Email:       "rumahpc.admin@gmail.com",
			Password:    "rumahpckuprofit",
			PhoneNumber: "08123456789",
			RoleID:      "d9ed6ac7-2a90-4a82-b2a9-c40fda4a5d93",
		},
		{
			Name:        "RumahPC Customer",
			Email:       "rumahpc.customer@gmail.com",
			Password:    "akucustomer",
			PhoneNumber: "08123456789",
			RoleID:      "54e26dde-9cc3-45b2-beb5-8b9c02d51eaa",
		},
		{
			Name:        "RumahPC Merchant",
			Email:       "rumahpc.merchant@gmail.com",
			Password:    "akumerchant",
			PhoneNumber: "08123456789",
			RoleID:      "9e1e7628-27c4-49b7-9204-a7b16a009edc",
		},
	}

	for _, user := range users {
		hashedPassword, salt, err := services.HashPassword(user.Password)

		if err != nil {
			fmt.Println(err)
			return
		}

		config.DB.FirstOrCreate(&user, models.User{
			ID:          uuid.New().String(),
			Name:        user.Name,
			Email:       user.Email,
			Password:    hashedPassword,
			Salt:        salt,
			PhoneNumber: user.PhoneNumber,
			RoleID:      user.RoleID,
		})
	}
}
