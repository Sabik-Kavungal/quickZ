package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"quickZ/handlers"
	"quickZ/middleware"
)

// AdminRoutes - Define all routes accessible by admin users
func AdminRoutes(r *gin.Engine, db *sql.DB) {
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware("admin"))
	{
		// Add a new product
		admin.POST("/add", handlers.AddProduct(db))

		// List all products
		admin.GET("/products", handlers.ListProducts(db))

		// Get a single product by ID
		admin.GET("/product/:id", handlers.GetProductByID(db))

		// Update an existing product by ID
		admin.PUT("/product/:id", handlers.UpdateProduct(db))

		// Delete a product by ID
		admin.DELETE("/product/:id", handlers.DeleteProduct(db))
	}
}
