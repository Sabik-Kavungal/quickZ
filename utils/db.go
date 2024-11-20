package utils

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:sabik123@tcp(127.0.0.1:3306)/quickz")
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	return db
}
