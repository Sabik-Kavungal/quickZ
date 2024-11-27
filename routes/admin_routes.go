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
		admin.GET("/products", handlers.ListProductsAndByCategory(db))

		// Get a single product by ID
		admin.GET("/product/:id", handlers.GetProductByID(db))

		// Update an existing product by ID
		admin.PUT("/product/:id", handlers.UpdateProduct(db))

		// Delete a product by ID
		admin.DELETE("/product/:id", handlers.DeleteProduct(db))

		// Get all users
		admin.GET("/users", handlers.GetAllUsers(db))

		// Add a new user
		admin.POST("/user", handlers.AddUser(db))

		// Update user by ID
		admin.PUT("/user/:id", handlers.UpdateUser(db))

		// Delete user by ID
		admin.DELETE("/user/:id", handlers.DeleteUser(db))

		// Category Routes
		// Add a new category
		admin.POST("/category", handlers.CreateCategory(db))

		// List all categories
		admin.GET("/categories", handlers.ListCategories(db))

		// Get a single category by ID
		admin.GET("/category/:id", handlers.GetCategoryByID(db))

		// Update an existing category by ID
		admin.PUT("/category/:id", handlers.UpdateCategory(db))

		// Delete a category by ID
		admin.DELETE("/category/:id", handlers.DeleteCategory(db))

		//chat with ai, ai chatbot
		admin.POST("/ai-chatbot", handlers.AiChatbot)

	}
}
