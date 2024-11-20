package main

import (
	"quickZ/handlers"
	"quickZ/middleware"
	"quickZ/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db := utils.ConnectDB()
	defer db.Close()

	r := gin.Default()

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://your-production-domain.com"}, // Replace with allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Authentication routes
	r.POST("/register", handlers.Register(db))
	r.POST("/login", handlers.Login(db))

	// Protected routes
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware("admin")) // Ensure only "admin" can manage products
	{
		admin.POST("/add", handlers.AddProduct(db))
	}

	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware("user")) // Ensure only "admin" can manage products
	{
		user.GET("/list", handlers.ListProducts(db))
	}

	r.Run(":8080")
}
