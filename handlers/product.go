package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"quickZ/models"
)

// AddProduct handles adding a new product to the database
func AddProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		// Parse the incoming JSON body
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Extract user ID from context (assuming context has userID set)
		createdBy := r.Context().Value("userID")
		if createdBy == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// PostgreSQL insert query with positional placeholders
		query := `INSERT INTO products (name, description, price, created_by) VALUES ($1, $2, $3, $4) RETURNING id`
		var productID int
		err := db.QueryRow(query, product.Name, product.Description, product.Price, createdBy).Scan(&productID)
		if err != nil {
			http.Error(w, "Could not add product", http.StatusInternalServerError)
			return
		}

		// Return the ID of the newly created product
		response := map[string]interface{}{
			"message":   "Product added successfully",
			"productID": productID,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// ListProducts handles listing all products from the database
func ListProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query to fetch products
		rows, err := db.Query("SELECT id, name, description, price, created_by FROM products")
		if err != nil {
			http.Error(w, "Could not fetch products", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var product models.Product
			// Scan the product fields
			if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CreatedBy); err != nil {
				http.Error(w, "Could not parse product", http.StatusInternalServerError)
				return
			}
			products = append(products, product)
		}

		// Check if any products were found
		if len(products) == 0 {
			http.Error(w, "No products found", http.StatusOK)
			return
		}
		
		// Return the list of products
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}