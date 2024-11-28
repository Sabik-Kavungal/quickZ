package handlers

import (
	"database/sql"
	"net/http"
	"quickZ/models"
	"strconv"
	"github.com/gin-gonic/gin"
)

func AddProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Validate required fields
		if product.CategoryID == 0 || product.ImageURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID and Image URL are required"})
			return
		}

		createdBy, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Insert product into the database with imageUrl
		query := `INSERT INTO products (name, description, price, created_by, category_id, image_url) 
		          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
		var productID int
		err := db.QueryRow(query, product.Name, product.Description, product.Price, createdBy, product.CategoryID, product.ImageURL).Scan(&productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product added successfully", "productID": productID})
	}
}
func ListProductsAndByCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryID := c.Query("category_id")

		var query string
		var rows *sql.Rows
		var err error

		if categoryID != "" {
			// Query for products filtered by category ID
			query = `SELECT p.id, p.name, p.description, p.price, p.created_by, 
			         p.category_id, p.image_url, c.name AS category_name
			         FROM products p
			         LEFT JOIN categories c ON p.category_id = c.id
			         WHERE p.category_id = $1`
			rows, err = db.Query(query, categoryID)
		} else {
			// Query for all products
			query = `SELECT p.id, p.name, p.description, p.price, p.created_by, 
			         p.category_id, p.image_url, c.name AS category_name
			         FROM products p
			         LEFT JOIN categories c ON p.category_id = c.id
			         ORDER BY c.name, p.id`
			rows, err = db.Query(query)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch products"})
			return
		}
		defer rows.Close()

		// Prepare the response structure
		type CategoryWiseProducts struct {
			Category string           `json:"category"`
			Products []models.Product `json:"products"`
		}

		var categoryMap = make(map[string][]models.Product)
		for rows.Next() {
			var product models.Product
			var categoryID sql.NullInt64  // Handle potential NULL values for category_id
			var imageURL sql.NullString   // Handle potential NULL values for image_url
			var categoryName sql.NullString // Handle potential NULL values for category_name

			// Scan the data
			if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price,
				&product.CreatedBy, &categoryID, &imageURL, &categoryName); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse product", "details": err.Error()})
				return
			}

			// Assign values from sql.NullInt64 and sql.NullString
			if categoryID.Valid {
				product.CategoryID = int(categoryID.Int64)
			} else {
				product.CategoryID = 0 // Default to 0 if NULL
			}

			if imageURL.Valid {
				product.ImageURL = imageURL.String
			} else {
				product.ImageURL = "" // Default to an empty string if NULL
			}

			category := categoryName.String
			if !categoryName.Valid {
				category = "Uncategorized" // Default category if NULL
			}

			// Add the product to the appropriate category group
			categoryMap[category] = append(categoryMap[category], product)
		}

		if len(categoryMap) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No products found"})
			return
		}

		// Convert the category map to a slice for response
		var groupedProducts []CategoryWiseProducts
		for category, products := range categoryMap {
			groupedProducts = append(groupedProducts, CategoryWiseProducts{
				Category: category,
				Products: products,
			})
		}

		// Return the grouped products
		c.JSON(http.StatusOK, groupedProducts)
	}
}



func GetProductByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the product ID from the URL parameter
		productID := c.Param("id")

		// SQL query to retrieve the product by its ID
		query := `SELECT p.id, p.name, p.description, p.price, p.imageUrl, p.created_by, p.category_id, 
		          c.name AS category_name
		          FROM products p
		          LEFT JOIN categories c ON p.category_id = c.id
		          WHERE p.id = $1`

		// Query the database for the product by ID
		row := db.QueryRow(query, productID)

		// Initialize the Product struct to hold the result
		var product models.Product
		var categoryName string

		// Scan the result into the Product struct and category name
		err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImageURL,
			&product.CreatedBy, &product.CategoryID, &categoryName)
		if err != nil {
			// If the product isn't found, return an error
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}
			// If there's another error (e.g., a database issue), return a 500
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch product"})
			return
		}

		// Set the category name in the product
		product.Category = categoryName

		// Return the product with category name
		c.JSON(http.StatusOK, product)
	}
}

// UpdateProduct updates the details of an existing product
func UpdateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve productID from the URL parameter
		productID := c.Param("id")

		// Initialize the Product struct to bind JSON data
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Convert the productID to an integer
		id, err := strconv.Atoi(productID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		// SQL query to update the product, including imageUrl
		query := `UPDATE products 
		          SET name = $1, description = $2, price = $3, category_id = $4, imageUrl = $5 
		          WHERE id = $6`

		// Execute the update query
		result, err := db.Exec(query, product.Name, product.Description, product.Price, product.CategoryID, product.ImageURL, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update product"})
			return
		}

		// Check if any rows were affected by the update
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking affected rows"})
			return
		}

		// If no rows were affected, return a not found error
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		// Return a success message
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
