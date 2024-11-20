package routes

import (
	"github.com/gin-gonic/gin"
	"quickZ/handlers"
	"quickZ/middleware"
	"database/sql"
)

// UserRoutes - Define all routes accessible by user
func UserRoutes(r *gin.Engine, db *sql.DB) {
	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware("user"))
	{
		user.GET("/list", handlers.ListProducts(db))
	}
}
