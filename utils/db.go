package utils

import (
	"database/sql"
	
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// ConnectDB establishes the connection to the PostgreSQL database
func ConnectDB() *sql.DB {
	// Use the provided connection string
	connStr := "postgresql://root:MBz1sFV9wu6XL6UDAVvRrYZY4UU6bKbq@dpg-csv3asjtq21c73eje9ng-a.singapore-postgres.render.com:5432/quick_ubvt?sslmode=require"

	// Log the connection string (for debugging purposes, remove in production)
	log.Println("Connecting to DB with connection string: ", connStr)

	// Open the connection to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	log.Println("Successfully connected to the PostgreSQL database")
	return db
}
