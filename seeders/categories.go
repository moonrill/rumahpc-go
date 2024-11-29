package seeders

import (
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
)

func SeedCategories() {
	categories := []models.Category{
		{Name: "Computer", Icon: "computer.png"},
		{Name: "Laptop", Icon: "laptop.png"},
		{Name: "PC Parts", Icon: "pc-parts.png"},
	}

	for _, category := range categories {
		config.DB.FirstOrCreate(&category, models.Category{Name: category.Name, Icon: category.Icon})
	}
}
