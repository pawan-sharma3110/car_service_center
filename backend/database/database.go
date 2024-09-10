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
	// Table Create Function
	createUserTable(db)

	createServiceTable(db)

	createAppoitmentsTable(db)
	return db
}
func createUserTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS users (
	 user_id UUID PRIMARY KEY,
    full_name VARCHAR(100),
    phone_no VARCHAR(15),
    email VARCHAR(100),
    password VARCHAR(255),
	role VARCHAR(15),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)

	}

}

func createServiceTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS services (
			  service_id UUID PRIMARY KEY,
				name VARCHAR(100),
				cost FLOAT,
				description TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	        )`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)

	}

}
func createAppoitmentsTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS appointments  (
	    appointment_id UUID PRIMARY KEY,
		user_id UUID REFERENCES users(user_id),
		services JSONB,  -- Store services as a JSON array
		date TIMESTAMP,  -- Store date in local time without timezone
		status VARCHAR(20) DEFAULT 'Unscheduled',
		total_cost FLOAT,
		created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)

	}

}
