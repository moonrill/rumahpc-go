package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/moonrill/rumahpc-api/internal/controllers"
	"github.com/moonrill/rumahpc-api/middleware"
	"github.com/moonrill/rumahpc-api/utils"
)

func SetupRoutes(router *gin.Engine) {
	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")
	{
		v1.GET("/image/*path", utils.ServeImage)

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
			subcategories.GET("/:slug", controllers.GetSubCategoryBySlug)
		}

		brands := v1.Group("/brands")
		{
			brands.GET("/", controllers.GetBrands)
			brands.GET("/:slug", controllers.GetBrandBySlug)
		}

		product := v1.Group("/products")
		{
			product.GET("/", controllers.GetProducts)
			product.GET("/category/:slug", controllers.GetProductsByCategorySlug)
			product.GET("/subcategory/:slug", controllers.GetProductsBySubCategorySlug)
			product.GET("/:slug", controllers.GetProduct)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.HeaderAuthMiddleware)
		{
			protected.GET("/profile", controllers.GetProfile)

			admin := protected.Group("/")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.POST("/categories", controllers.CreateCategory)
				admin.POST("/categories/upload", controllers.UploadCategoryIcon)
				admin.PUT("/categories/:id", controllers.UpdateCategory)
				admin.DELETE("/categories/:id", controllers.DeleteCategory)

				admin.POST("/subcategories", controllers.CreateSubCategory)
				admin.POST("/subcategories/upload", controllers.UploadSubCategoryIcon)
				admin.PUT("/subcategories/:id", controllers.UpdateSubCategory)
				admin.DELETE("/subcategories/:id", controllers.DeleteSubCategory)

				admin.POST("/brands", controllers.CreateBrand)
				admin.POST("/brands/upload", controllers.UploadBrandIcon)
				admin.PUT("/brands/:id", controllers.UpdateBrand)
				admin.DELETE("/brands/:id", controllers.DeleteBrand)
			}

			merchant := protected.Group("/")
			merchant.Use(middleware.RoleMiddleware("merchant"))
			{
				merchant.POST("/products", controllers.CreateProduct)
				merchant.POST("/products/upload", controllers.UploadProductImage)
				merchant.POST("/products/upload/multiple", controllers.UploadMultipleProductImages)
				merchant.PUT("/products/status/:id", controllers.ToggleProductStatus)
				merchant.PUT("/products/:id", controllers.UpdateProduct)
			}

			customer := protected.Group("/")
			customer.Use(middleware.RoleMiddleware("customer"))
			{
				customer.GET("/carts", controllers.GetCart)
				customer.GET("/carts/grouped", controllers.GetGroupedCart)
				customer.POST("/carts", controllers.AddToCart)
				customer.PUT("/carts/:id", controllers.UpdateCart)
				customer.DELETE("/carts", controllers.RemoveFromCart)

				customer.GET("/shipping/rates/cart", controllers.GetCartCouriersRates)
				customer.GET("/shipping/rates/buy-now", controllers.GetBuyNowCouriersRates)

				customer.POST("/orders/buy-now", controllers.CreateBuyNowOrder)
				customer.POST("/orders/cart", controllers.CreateCartOrder)

			}

			customerMerchant := protected.Group("/")
			customerMerchant.Use(middleware.RoleMiddleware("customer", "merchant"))
			{
				customerMerchant.GET("/addresses", controllers.GetAddress)
				customerMerchant.GET("/addresses/:id", controllers.GetAddressById)
				customerMerchant.POST("/addresses", controllers.CreateAddress)
				customerMerchant.PUT("/addresses/:id", controllers.UpdateAddress)
				customerMerchant.DELETE("/addresses/:id", controllers.DeleteAddress)
			}
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "API is up and running",
		})
	})
}
