package routes

import (
	"bicycle-store/internal/controllers"
	"bicycle-store/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine) {
	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		authController := controllers.NewAuthController()
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.GET("/me", middleware.AuthMiddleware(), authController.GetMe)
		}

		// Category routes
		categoryController := controllers.NewCategoryController()
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryController.GetAll)
			categories.GET("/:id", categoryController.GetByID)
			// Admin only
			categories.POST("", middleware.AuthMiddleware(), middleware.AdminMiddleware(), categoryController.Create)
			categories.PUT("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), categoryController.Update)
			categories.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), categoryController.Delete)
		}

		// Bicycle routes
		bicycleController := controllers.NewBicycleController()
		bicycles := v1.Group("/bicycles")
		{
			bicycles.GET("", bicycleController.GetAll)
			bicycles.GET("/:id", bicycleController.GetByID)
			// Admin only
			bicycles.POST("", middleware.AuthMiddleware(), middleware.AdminMiddleware(), bicycleController.Create)
			bicycles.PUT("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), bicycleController.Update)
			bicycles.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), bicycleController.Delete)
			bicycles.PATCH("/:id/stock", middleware.AuthMiddleware(), middleware.AdminMiddleware(), bicycleController.UpdateStock)
			// Customer - add review
			bicycles.POST("/:id/reviews", middleware.AuthMiddleware(), bicycleController.AddReview)
		}

		// Order routes
		orderController := controllers.NewOrderController()
		orders := v1.Group("/orders")
		orders.Use(middleware.AuthMiddleware())
		{
			orders.GET("/my", orderController.GetMyOrders)
			orders.POST("", orderController.Create)
			orders.GET("/:id", orderController.GetByID)
			// Admin only
			orders.GET("", middleware.AdminMiddleware(), orderController.GetAll)
			orders.PATCH("/:id/status", middleware.AdminMiddleware(), orderController.UpdateStatus)
		}

		// Customer routes
		customerController := controllers.NewCustomerController()
		customers := v1.Group("/customers")
		customers.Use(middleware.AuthMiddleware())
		{
			customers.PUT("/profile", customerController.UpdateProfile)
			customers.POST("/addresses", customerController.AddAddress)
			customers.DELETE("/addresses/:type", customerController.RemoveAddress)
			// Admin only
			customers.GET("", middleware.AdminMiddleware(), customerController.GetAll)
			customers.GET("/:id", middleware.AdminMiddleware(), customerController.GetByID)
		}

		// Report routes (Admin only)
		reportController := controllers.NewReportController()
		reports := v1.Group("/reports")
		reports.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			reports.GET("/sales-by-category", reportController.GetSalesByCategory)
			reports.GET("/top-selling", reportController.GetTopSellingBicycles)
		}
	}
}
