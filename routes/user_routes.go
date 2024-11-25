package routes

import (
	"database/sql"
	"quickZ/handlers"
	"quickZ/middleware"
	"github.com/gorilla/mux"
)

// UserRoutes - Define all routes accessible by user
func UserRoutes(r *mux.Router, db *sql.DB) {
	// Create a sub-router for user routes
	user := r.PathPrefix("/user").Subrouter()

	// Apply middleware to user routes (auth for user)
	user.Use(middleware.AuthMiddleware("user"))

	// Define the route for listing products
	user.HandleFunc("/list", handlers.ListProducts(db)).Methods("GET")
}
