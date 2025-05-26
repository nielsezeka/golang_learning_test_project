package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	connStr := "host=localhost port=5432 user=postgres dbname=test_db sslmode=disable search_path=public"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Log the current database and user for debugging
	var currentDB, currentUser string
	_ = DB.QueryRow("SELECT current_database(), current_user").Scan(&currentDB, &currentUser)
	fmt.Printf("Connected to PostgreSQL! Current DB: %s, User: %s\n", currentDB, currentUser)
}
