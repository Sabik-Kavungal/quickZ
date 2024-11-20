package routes

import (
	"github.com/gin-gonic/gin"
	"quickZ/handlers"
	"quickZ/middleware"
	"database/sql"
)

// AdminRoutes - Define all routes accessible by admin users
func AdminRoutes(r *gin.Engine, db *sql.DB) {
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware("admin"))
	{
		admin.POST("/add", handlers.AddProduct(db))
		admin.GET("/list", handlers.ListProducts(db)) // For demonstration, you could have other admin functionalities
	}
}
