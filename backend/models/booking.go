package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`    // Reference to the user who made the appointment
	Services  []Service `json:"services"`   // List of services chosen for the appointment
	Date      time.Time `json:"date"`       // Date and time of the appointment
	Status    string    `json:"status"`     // Status of the appointment (Scheduled, Completed, Canceled)
	TotalCost float64   `json:"total_cost"` // Total cost of the selected services
	CreatedOn time.Time `json:"created_on"` // Timestamp when the appointment was created
}

// Helper function to convert services to JSON
func ServicesToJson(services []Service) string {
	servicesJson, _ := json.Marshal(services)
	return string(servicesJson)
}
func (a Appointment) InsertAppoitment(db *sql.DB) error {
	query := `
        INSERT INTO appointments (id, user_id, services, date, status, total_cost, created_on)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(query,
		a.Id,
		a.UserId,
		ServicesToJson(a.Services), // Convert services to JSON format
		a.Date,
		a.Status,
		a.TotalCost,
		a.CreatedOn,
	)
	if err != nil {

		return fmt.Errorf("failed to create appointment: %w", err)
	}
	return nil
}
