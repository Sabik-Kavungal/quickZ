package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"quickZ/models"
	"quickZ/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
	"golang.org/x/crypto/bcrypt"
)

// Register handler for user registration
func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Hash the password using bcrypt
		//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		// 	return
		// }

		// Insert user into the PostgreSQL database
		// 		_, err = db.Exec("INSERT INTO users (username, password, type) VALUES ($1, $2, $3)", user.Username, user.Password, user.Type)
		// 		if err != nil {
		// 			// Log the actual error message for debugging
		// 			log.Printf("Error saving user: %v", err)
		// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
		// 			return
		// 		}

		// 		// Return success message
		// 		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
		// 	}
		// }
		// No hashing for password, use it as is
		_, err := db.Exec("INSERT INTO users (username, password, type) VALUES ($1, $2, $3)", user.Username, user.Password, user.Type)
		if err != nil {
			log.Printf("Error saving user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User added successfully"})
	}
}

// Login handler for user login
func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData models.User
		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var user models.User
		// Query the database for the user by username
		err := db.QueryRow("SELECT id, password, type FROM users WHERE username = $1", loginData.Username).
			Scan(&user.ID, &user.Password, &user.Type)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// Compare the password with the stored hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
			return
		}

		// Generate a JWT token
		token, _ := utils.GenerateToken(user.ID, user.Type)

		// Return the token and user data
		c.JSON(http.StatusOK, gin.H{"token": token, "data": user})
	}
}

// Get all users (Admin only)
func GetAllUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, username,password, type FROM users")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Username, &user.Type); err != nil {
				log.Printf("Error scanning user: %v", err)
				continue
			}
			users = append(users, user)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading rows"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": users})
	}
}

// Update user (Admin only)
func UpdateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Optional: Hash password if it is updated
		if user.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
				return
			}
			user.Password = string(hashedPassword)
		}

		// Update user information
		_, err := db.Exec("UPDATE users SET username = $1, password = $2, type = $3 WHERE id = $4",
			user.Username, user.Password, user.Type, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	}
}

// Add user (Admin only)
func AddUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		// Insert user into the database
		_, err = db.Exec("INSERT INTO users (username, password, type) VALUES ($1, $2, $3)", user.Username, hashedPassword, user.Type)
		if err != nil {
			log.Printf("Error saving user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User added successfully"})
	}
}

// Delete user (Admin only)
func DeleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Delete user from database
		_, err := db.Exec("DELETE FROM users WHERE id = $1", c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
