package seeders

import (
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
)

func SeedCategories() {
	categories := []models.Category{
		{Name: "Computer", Slug: "computer", Icon: "computer.png"},
		{Name: "Laptop", Slug: "laptop", Icon: "laptop.png"},
		{Name: "PC Parts", Slug: "pc-parts", Icon: "pc-parts.png"},
	}

	for _, category := range categories {
		config.DB.FirstOrCreate(&category, models.Category{Name: category.Name, Icon: category.Icon})
	}
}
