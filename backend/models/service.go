package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ServiceID   uuid.UUID `json:"service_id"`
	Name        string    `json:"name"`
	Cost        float32   `json:"cost"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s Service) InsertService(db *sql.DB) error {
	_, err := db.Exec(`INSERT INTO services(service_id, name, cost, description, created_at) VALUES($1, $2, $3, $4, $5)`, s.ServiceID, s.Name, s.Cost, s.Description, s.CreatedAt)
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("failed to insert service into database: %w", err)
	}

	return nil
}

func AllService(db *sql.DB) ([]Service, error) {
	// Query to get all services
	rows, err := db.Query(`SELECT service_id, name, cost, description, created_at FROM services`)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve services: %w", err)
	}
	defer rows.Close()

	var services []Service

	// Iterate over the result set
	for rows.Next() {
		var s Service
		// Scan each row into the Service struct
		err := rows.Scan(&s.ServiceID, &s.Name, &s.Cost, &s.Description, &s.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service: %w", err)
		}
		// Append each service to the slice
		services = append(services, s)
	}

	// Check for errors after iterating through the rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over services: %w", err)
	}

	return services, nil
}

func DeleteService(id uuid.UUID, db *sql.DB) error {
	row, err := db.Exec(`DELETE FROM services WHERE service_id = $1`, id)
	if err != nil {

		return err
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Service not found or maybe already deleted")
	}
	return nil
}

func DeleteAll(db *sql.DB) (err error) {
	row, err := db.Exec(`DELETE FROM services`)
	if err != nil {
		return fmt.Errorf("failed to delete all service: %w", err)
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Service not found or maybe already deleted")
	}
	return nil
}
