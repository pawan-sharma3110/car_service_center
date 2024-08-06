package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func DbIn() (db *sql.DB) {
	// Connection string
	conStr := "host=localhost port=5432 user=postgres password=Pawan@2003 dbname=car_service sslmode=disable"

	// Open the database
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
		return nil
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("error pinging database: %v", err)
	}

	return db
}
