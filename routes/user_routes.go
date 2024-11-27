package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"quickZ/handlers"
	"quickZ/middleware"
)

// UserRoutes - Define all routes accessible by user
func UserRoutes(r *gin.Engine, db *sql.DB) {
	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware("user"))
	{
		user.GET("/list", handlers.ListProductsAndByCategory(db))
	}
}
