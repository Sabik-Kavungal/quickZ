package handlers

import (
	"database/sql"
	"net/http"
	"quickZ/models"

	"github.com/gin-gonic/gin"
)

// AddProduct adds a new product to the database
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

		query := `INSERT INTO products (name, description, price, created_by) VALUES ($1, $2, $3, $4) RETURNING id`
		var productID int
		err := db.QueryRow(query, product.Name, product.Description, product.Price, createdBy).Scan(&productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product added successfully", "productID": productID})
	}
}

// ListProducts retrieves all products from the database
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

		if len(products) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No products found"})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}

// GetProductByID retrieves a single product by its ID
func GetProductByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")
		var product models.Product
		query := "SELECT id, name, description, price, created_by FROM products WHERE id = $1"
		err := db.QueryRow(query, productID).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CreatedBy)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch product"})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

// UpdateProduct updates the details of an existing product
func UpdateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		query := `UPDATE products SET name = $1, description = $2, price = $3 WHERE id = $4`
		result, err := db.Exec(query, product.Name, product.Description, product.Price, productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update product"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
	}
}

// DeleteProduct removes a product from the database
func DeleteProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")
		query := `DELETE FROM products WHERE id = $1`
		result, err := db.Exec(query, productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete product"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	}
}
