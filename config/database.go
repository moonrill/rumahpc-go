package config

import (
	"fmt"
	"log"
	"os"

	"github.com/moonrill/rumahpc-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	// Load environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Format the database connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta", host, port, user, password, dbname)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(
		&models.Category{},
		&models.User{},
		&models.SubCategory{},
		&models.Brand{},
		&models.Role{},
		&models.Address{},
		&models.Product{},
		&models.ProductImages{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.Cart{},
		&models.CartItem{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = db

	log.Println("Connected to database successfully")
}

func GetDB() *gorm.DB {
	return DB
}
