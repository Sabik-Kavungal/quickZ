package models

// Category represents the structure of a category in the database.
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
