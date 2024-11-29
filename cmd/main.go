package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/routes"
	"github.com/moonrill/rumahpc-api/seeders"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	// Initialize the database
	config.InitDatabase()

	// Router setup
	router := gin.Default()

	// Load routes
	routes.SetupRoutes(router)

	// Load seeders
	seeders.RunSeeders()

	// Get Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Running the server
	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
