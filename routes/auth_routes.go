package routes

import (
	"database/sql"
	"quickZ/handlers"
	"github.com/gorilla/mux"
)

// AuthRoutes - Define routes related to authentication
func AuthRoutes(r *mux.Router, db *sql.DB) {
	// POST /register route for user registration
	r.HandleFunc("/register", handlers.Register(db)).Methods("POST")

	// POST /login route for user login
	r.HandleFunc("/login", handlers.Login(db)).Methods("POST")
}
