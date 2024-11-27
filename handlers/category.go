package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"quickZ/models"
)

// CreateCategory adds a new category to the database
func CreateCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category models.Category
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Insert category into the database
		query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
		var categoryID int
		err := db.QueryRow(query, category.Name).Scan(&categoryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add category"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category created successfully", "categoryID": categoryID})
	}
}

// ListCategories retrieves all categories from the database
func ListCategories(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `SELECT id, name FROM categories`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch categories"})
			return
		}
		defer rows.Close()

		var categories []models.Category
		for rows.Next() {
			var category models.Category
			if err := rows.Scan(&category.ID, &category.Name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse category"})
				return
			}
			categories = append(categories, category)
		}

		if len(categories) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No categories found"})
			return
		}

		c.JSON(http.StatusOK, categories)
	}
}

// GetCategoryByID retrieves a category by its ID
func GetCategoryByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryID := c.Param("id")
		var category models.Category

		query := `SELECT id, name FROM categories WHERE id = $1`
		err := db.QueryRow(query, categoryID).Scan(&category.ID, &category.Name)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch category"})
			return
		}

		c.JSON(http.StatusOK, category)
	}
}

// UpdateCategory updates the details of an existing category
func UpdateCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryID := c.Param("id")
		var category models.Category
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		query := `UPDATE categories SET name = $1 WHERE id = $2`
		result, err := db.Exec(query, category.Name, categoryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update category"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
	}
}

// DeleteCategory removes a category from the database
func DeleteCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryID := c.Param("id")
		query := `DELETE FROM categories WHERE id = $1`
		result, err := db.Exec(query, categoryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete category"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
	}
}
