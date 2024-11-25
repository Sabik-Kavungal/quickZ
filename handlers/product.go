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

		// Extract user ID from the context
		createdBy, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// PostgreSQL insert query with positional placeholders
		query := `INSERT INTO products (name, description, price, created_by) VALUES ($1, $2, $3, $4) RETURNING id`
		var productID int
		err := db.QueryRow(query, product.Name, product.Description, product.Price, createdBy).Scan(&productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add product"})
			return
		}

		// Return the ID of the newly created product
		c.JSON(http.StatusOK, gin.H{"message": "Product added successfully", "productID": productID})
	}
}

func ListProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Query to fetch products
		rows, err := db.Query("SELECT id, name, description, price, created_by FROM products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch products"})
			return
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var product models.Product
			// Scan the product fields
			if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CreatedBy); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse product"})
				return
			}
			products = append(products, product)
		}

		// Check if any products were found
		if len(products) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No products found &&"})
			return
		}

		// Return the list of products
		c.JSON(http.StatusOK, products)
	}
}
