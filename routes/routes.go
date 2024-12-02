package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/controllers"
	"github.com/moonrill/rumahpc-api/middleware"
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
			categories.GET("/:slug", controllers.GetCategoryBySlug)
		}

		subcategories := v1.Group("/subcategories")
		{
			subcategories.GET("/", controllers.GetSubCategories)
			subcategories.GET("/:slug", controllers.GetSubCategoriesBySlug)
		}

		brands := v1.Group("/brands")
		{
			brands.GET("/", controllers.GetBrands)
		}

		// Protected routes
		protected := v1.Group("/")
		// TODO: choose the auth middleware
		// protected.Use(middleware.CookiesAuthMiddleware)
		protected.Use(middleware.HeaderAuthMiddleware)
		{
			protected.GET("/profile", controllers.GetProfile)

			protected.POST("/categories", controllers.CreateCategory)
			protected.PUT("/categories/:id", controllers.UpdateCategory)
			protected.DELETE("/categories/:id", controllers.DeleteCategory)

			protected.POST("/subcategories", controllers.CreateSubCategory)
			protected.PUT("/subcategories/:id", controllers.UpdateSubCategory)
			protected.DELETE("/subcategories/:id", controllers.DeleteSubCategory)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "API is up and running",
		})
	})
}
