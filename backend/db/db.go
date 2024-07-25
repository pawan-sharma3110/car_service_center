package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func DbIn() (db *sql.DB, err error) {
	// Connection string
	conStr := "host=localhost port=5432 user=postgres password=Pawan@2003 dbname=car_service sslmode=disable"

	// Open the database
	db, err = sql.Open("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return db, nil
}
