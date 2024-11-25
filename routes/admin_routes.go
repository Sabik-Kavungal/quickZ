package routes

import (
	"database/sql"
	"quickZ/handlers"
	"quickZ/middleware"
	"github.com/gorilla/mux"
)

// AdminRoutes - Define all routes accessible by admin users
func AdminRoutes(r *mux.Router, db *sql.DB) {
	// Create a new sub-router for admin routes
	admin := r.PathPrefix("/admin").Subrouter()

	// Apply middleware to the admin routes (auth for admin)
	admin.Use(middleware.AuthMiddleware("admin"))

	// Define the routes
	admin.HandleFunc("/add", handlers.AddProduct(db)).Methods("POST")
	admin.HandleFunc("/list", handlers.ListProducts(db)).Methods("GET") // For demonstration, you could have other admin functionalities
}
