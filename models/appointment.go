package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Appointment struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`    // Reference to the user who made the appointment
	Services  []Service `json:"services"`   // List of services chosen for the appointment
	DateTime  string    `json:"date_time"`  // Date and time of the appointment
	Status    string    `json:"status"`     // Status of the appointment (Scheduled, Completed, Canceled)
	TotalCost float64   `json:"total_cost"` // Total cost of the selected services
	CreatedOn string    `json:"created_on"` // Timestamp when the appointment was created
}

// Helper function to convert services to JSON
func ServicesToJson(services []Service) string {
	servicesJson, _ := json.Marshal(services)
	return string(servicesJson)
}

func AllAppointment(db *sql.DB) (*[]Appointment, error) {
	// Query to fetch all appointments
	query := `SELECT id, user_id, services, date, status, total_cost, created_on FROM appointments`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch appointments %w", err)
	}
	defer rows.Close()

	var appointments []Appointment

	for rows.Next() {
		var appointment Appointment
		var servicesJson string

		// Scan the appointment details
		err = rows.Scan(
			&appointment.ID,
			&appointment.UserID,
			&servicesJson, // This stores services as JSON string from DB
			&appointment.DateTime,
			&appointment.Status,
			&appointment.TotalCost,
			&appointment.CreatedOn,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appointments %w", err)
		}

		// Parse services from JSON into Go struct
		err = json.Unmarshal([]byte(servicesJson), &appointment.Services)
		if err != nil {
			return nil, fmt.Errorf("failed to parse services %w", err)
		}
		// Append appointment to the list
		fmt.Println(appointment.DateTime)
		appointments = append(appointments, appointment)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during iteration %w", err)
	}
	return &appointments, nil
}

func (a Appointment) UpdateStatus(db *sql.DB) error {

	// Prepare the SQL query to update the appointment status and accepted date
	query := `UPDATE appointments SET status = $1, date = $2 WHERE id = $3`
	result, err := db.Exec(query, a.Status, a.DateTime, a.ID)
	if err != nil {

		return fmt.Errorf("failed to update appointment %w", err)
	}

	// Check if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {

		return fmt.Errorf("error checking update result %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no appointment found with the given ID %w", err)
	}

	return nil
}
