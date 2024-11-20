package handlers

import (
	"database/sql"
	"net/http"
	"quickZ/models"

	"github.com/gin-gonic/gin"
)
func AddProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		createdBy, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		query := `INSERT INTO products (name, description, price, created_by) VALUES (?, ?, ?, ?)`
		_, err := db.Exec(query, product.Name, product.Description, product.Price, createdBy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product added successfully"})
	}
}

func ListProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, description, price, created_by FROM products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch products"})
			return
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var product models.Product
			if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CreatedBy); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse product"})
				return
			}
			products = append(products, product)
		}

		c.JSON(http.StatusOK, products)
	}
}