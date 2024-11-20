package handlers

import (
	"database/sql"
	"net/http"
	"quickZ/models"
	"quickZ/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		_, err := db.Exec("INSERT INTO users (username, password, type) VALUES (?, ?, ?)", user.Username, hashedPassword, user.Type)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	}
}

func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData models.User
		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var user models.User
		err := db.QueryRow("SELECT id, password, type FROM users WHERE username = ?", loginData.Username).
			Scan(&user.ID, &user.Password, &user.Type)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
			return
		}

		token, _ := utils.GenerateToken(user.ID, user.Type)
		c.JSON(http.StatusOK, gin.H{"token": token, "data": user})
	}
}
