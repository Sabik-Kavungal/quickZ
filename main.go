package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"quickZ/routes"
	"quickZ/utils"
)

func main() {
	// Initialize the database connection
	db := utils.ConnectDB()
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()

	// Set up CORS
	// Apply the CORS middleware globally
	r.Use(utils.EnableCORS())

	// Define Routes
	routes.AuthRoutes(r, db)  // Authentication routes (Register, Login)
	routes.AdminRoutes(r, db) // Admin routes (Product management)
	routes.UserRoutes(r, db)  // User routes (List products)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "2000" // Default port for local development
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
