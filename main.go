package main

import (
	"github.com/gin-gonic/gin"
	"log"
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
	r.Use(utils.SetupCORS())

	// Define Routes
	routes.AuthRoutes(r, db)  // Authentication routes (Register, Login)
	routes.AdminRoutes(r, db) // Admin routes (Product management)
	routes.UserRoutes(r, db)  // User routes (List products)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
