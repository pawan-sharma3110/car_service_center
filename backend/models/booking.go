package models

import (
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	Id          uuid.UUID   `json:"id"`
	UserId      uuid.UUID   `json:"user_id"`         // Reference to the user who made the appointment
	Services    []Service   `json:"services"`        // List of services chosen for the appointment
	Date        time.Time   `json:"date"`            // Date and time of the appointment
	Status      string      `json:"status"`          // Status of the appointment (e.g., Scheduled, Completed, Canceled)
	TotalCost   float64     `json:"total_cost"`      // Total cost of the selected services
	CreatedOn   time.Time   `json:"created_on"`      // Timestamp when the appointment was created
	
}

