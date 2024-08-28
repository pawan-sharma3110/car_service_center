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

	query := `CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		full_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		phone_no VARCHAR(20) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL
	)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
