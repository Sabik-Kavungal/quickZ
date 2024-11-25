package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"quickZ/models"
	"quickZ/utils"

	"golang.org/x/crypto/bcrypt"
)

// Register handler for user registration
func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		// Decode the incoming JSON payload into the 'user' struct
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Hash the password using bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		// Insert user into the PostgreSQL database
		_, err = db.Exec("INSERT INTO users (username, password, type) VALUES ($1, $2, $3)", user.Username, hashedPassword, user.Type)
		if err != nil {
			// Log the actual error message for debugging
			log.Printf("Error saving user: %v", err)
			http.Error(w, "Error saving user", http.StatusInternalServerError)
			return
		}

		// Return success message
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
	}
}

// Login handler for user login
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginData models.User
		// Decode the incoming JSON payload into the 'loginData' struct
		err := json.NewDecoder(r.Body).Decode(&loginData)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		var user models.User
		// Query the database for the user by username
		err = db.QueryRow("SELECT id, password, type FROM users WHERE username = $1", loginData.Username).
			Scan(&user.ID, &user.Password, &user.Type)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		// Compare the password with the stored hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
			http.Error(w, "Incorrect password", http.StatusUnauthorized)
			return
		}

		// Generate a JWT token
		token, _ := utils.GenerateToken(user.ID, user.Type)

		// Return the token and user data
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"token": token,
			"data":  user,
		}
		json.NewEncoder(w).Encode(response)
	}
}
