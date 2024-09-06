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
		return
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Appointment created successfully"})
}

func GetAllAppointmentsHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	appointments, err := models.AllAppointment(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Send the appointments as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
}

// Appointment Status Update for admin
func AppointmentStatusUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	var appointment models.Appointment
	var err error
	appointment.Id, err = uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}

	// Parse the request body

	err = json.NewDecoder(r.Body).Decode(&appointment)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ensure status and accepted_date are provided
	if appointment.Status == "" || appointment.Date.IsZero() {
		http.Error(w, "Missing status or accepted date", http.StatusBadRequest)
		return
	}
	err = appointment.UpdateStatus(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Appointment updated successfully"))
}
