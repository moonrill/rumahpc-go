package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/controllers"
)

func SetupRoutes(router *gin.Engine) {
	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controllers.SignUp)
			auth.POST("/login", controllers.SignIn)
		}

		categories := v1.Group("/categories")
		{
			categories.GET("/", controllers.GetCategories)
			categories.GET("/:id", controllers.GetCategoryByID)
			categories.POST("/", controllers.CreateCategory)
			categories.PUT("/:id", controllers.UpdateCategory)
			categories.DELETE("/:id", controllers.DeleteCategory)
		}

		subcategories := v1.Group("/subcategories")
		{
			subcategories.GET("/:id", controllers.GetSubCategoriesByCategoryID)
			subcategories.POST("/", controllers.CreateSubCategory)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "API is up and running",
		})
	})
}
