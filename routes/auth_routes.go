package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"quickZ/handlers"
)

// AuthRoutes - Define routes related to authentication
func AuthRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/register", handlers.Register(db))
	r.POST("/login", handlers.Login(db))
}
