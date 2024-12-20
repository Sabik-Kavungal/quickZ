package models

// Product represents the structure of a product in the database.
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CreatedBy   int     `json:"created_by"`
	CategoryID  int     `json:"category_id"`
	Category    string  `json:"category"`   // Store only the category name
	ImageURL    string  `json:"image_url"`  // URL of the product's image
}
