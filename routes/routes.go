package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/controllers"
)

func SetupRoutes(router *gin.Engine) {
	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Route group for the API
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		// auth := v1.Group("/auth")
		// {

		// }

		categories := v1.Group("/categories")
		{
			categories.GET("/", controllers.GetCategories)
			categories.GET("/:id", controllers.GetCategoryByID)
			categories.POST("/", controllers.CreateCategory)
			categories.PUT("/:id", controllers.UpdateCategory)
			categories.DELETE("/:id", controllers.DeleteCategory)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "API is up and running",
		})
	})
}
