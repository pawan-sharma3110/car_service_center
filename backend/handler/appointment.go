package handler

import (
	"car_service/models"
	"car_service/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Create a new appointment

func CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var appt models.Appointment
	err := json.NewDecoder(r.Body).Decode(&appt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Insert appointment into the DB
	appt.ID = uuid.New()
	appt.CreatedOn = time.Now()
	appt.Status = "Unscheduled" // Default status

	// Serialize services into JSON
	servicesJSON, err := json.Marshal(appt.Services)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query := `
        INSERT INTO appointments (appointment_id, user_id, services, date, status, total_cost, created_on)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.Exec(query, appt.ID, appt.UserID, servicesJSON, appt.Date, appt.Status, appt.TotalCost, appt.CreatedOn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appt)
}

// Get all appointments with user details
func GetAppointments(w http.ResponseWriter, r *http.Request) {
	// Define the SQL query to fetch appointments
	query := `
	 SELECT a.appointment_id, u.full_name, u.phone_no, a.services, a.total_cost, a.date, a.status
	 FROM appointments a
	 JOIN users u ON a.user_id = u.user_id  -- Change this to the correct column name if needed
	 `

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error executing query: %v", err) // Log the error for debugging
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Prepare a slice to hold appointment data
	appointments := []map[string]interface{}{}

	// Iterate over rows
	for rows.Next() {
		var id, fullName, phoneNo, services, status string
		var totalCost float32
		var date sql.NullTime

		// Scan the row into variables
		if err := rows.Scan(&id, &fullName, &phoneNo, &services, &totalCost, &date, &status); err != nil {
			log.Printf("Error scanning row: %v", err) // Log the error for debugging
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Append data to the slice
		appointments = append(appointments, map[string]interface{}{
			"id":         id,
			"full_name":  fullName,
			"phone_no":   phoneNo,
			"services":   services,
			"total_cost": totalCost,
			"date":       utils.FormatDate(date), // Handle date formatting
			"status":     status,
		})
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err) // Log the error for debugging
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set response header and encode JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(appointments); err != nil {
		log.Printf("Error encoding JSON: %v", err) // Log the error for debugging
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Appointment Status Update for admin

func UpdateAppointment(w http.ResponseWriter, r *http.Request) {

	// Extract appointment ID from URL path
	appointmentID := r.PathValue("id")
	if appointmentID == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	print(appointmentID)
	var appt models.Appointment
	err := json.NewDecoder(r.Body).Decode(&appt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
        UPDATE appointments 
        SET date = $1, status = $2
        WHERE appointment_id = $3`
	_, err = db.Exec(query, appt.Date, appt.Status, appointmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Appointment updated successfully"))
}

// Delete appointment by ID
func DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	appointmentID := r.URL.Query().Get("id")
	_, err := db.Exec("DELETE FROM appointments WHERE appointment_id = $1", appointmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Appointment deleted successfully"))
}
