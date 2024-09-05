package handler

import (
	"car_service/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CreateAppointment(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the incoming request body into an Appointment struct
	var appointment models.Appointment
	var err error
	appointment.UserId, err = uuid.Parse(appointment.UserId.String())
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&appointment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// "Invalid request payload"
	// Assign a new UUID for the appointment and current time for created_on
	appointment.Id = uuid.New()
	appointment.CreatedOn = time.Now()

	// Calculate the total cost of the services
	var totalCost float64
	for _, service := range appointment.Services {
		totalCost += float64(service.Cost)
	}
	appointment.TotalCost = totalCost

	// Insert appointment into the database
	err = appointment.InsertAppoitment(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Appointment created successfully"})
}
