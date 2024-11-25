package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"quickZ/routes"
	"quickZ/utils"
)

func main() {
	// Initialize the database connection
	db := utils.ConnectDB()
	defer db.Close()

	// Initialize Gorilla Mux router
	r := mux.NewRouter()

	// Set up CORS middleware globally
	r.Use(utils.EnableCORS)

	// Define Routes
	routes.AuthRoutes(r, db)  // Authentication routes (Register, Login)
	routes.AdminRoutes(r, db) // Admin routes (Product management)
	routes.UserRoutes(r, db)  // User routes (List products)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "2000" // Default port for local development
	}

	// Start the server with the Gorilla Mux router
	log.Printf("Server running on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
