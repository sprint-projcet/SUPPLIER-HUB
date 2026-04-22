package routes

import (
	"supplierhub-backend/controllers"
	"supplierhub-backend/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes mengatur pemetaan semua rute API dalam aplikasi Gin
func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// Auth Routes (Public)
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", controllers.Register)
		authGroup.POST("/login", controllers.Login)
	}

	// Semua rute di bawah ini wajib melampirkan JWT token
	api.Use(middlewares.RequireAuth())

	// UMKM (User) Routes
	userGroup := api.Group("/user")
	userGroup.Use(middlewares.RequireRole("user"))
	{
		userGroup.GET("/stats", controllers.GetUserStats)
		userGroup.GET("/orders", controllers.GetUserOrders)
		// Produk katalog (umkm viewing products)
		userGroup.GET("/products", controllers.GetProducts)
		userGroup.POST("/orders", controllers.CreateOrder)
	}

	// Supplier Routes
	supplierGroup := api.Group("/supplier")
	supplierGroup.Use(middlewares.RequireRole("supplier"))
	{
		supplierGroup.GET("/stats", controllers.GetSupplierStats)
		supplierGroup.GET("/products", controllers.GetSupplierProducts)
		// CRUD Produk biasanya akan ada Endpoint POST/PUT/DELETE juga di sini
		supplierGroup.GET("/orders", controllers.GetSupplierOrders)
		supplierGroup.PUT("/orders/:id", controllers.UpdateOrderStatus)
	}

	// Admin Routes
	adminGroup := api.Group("/admin")
	adminGroup.Use(middlewares.RequireRole("admin"))
	{
		adminGroup.GET("/stats", controllers.GetAdminStats)
		adminGroup.GET("/suppliers", controllers.GetAdminSuppliers)
		adminGroup.PUT("/suppliers/:id/verify", controllers.VerifySupplier)
		adminGroup.GET("/logs", controllers.GetAdminLogs)
	}
}
