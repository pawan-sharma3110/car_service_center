package handler

import (
	"car_service/models"
	"car_service/utils"
	"encoding/json"
	"log"

	"net/http"

	"github.com/google/uuid"
)

// Create a new appointment

func CreateAppointment(w http.ResponseWriter, r *http.Request) {
	// Fetch the user ID from the request (e.g., cookies or headers)
	userId, err := utils.GetUserIDFromCookies(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Parse the incoming request body to the Appointment struct
	var appt models.Appointment
	err = json.NewDecoder(r.Body).Decode(&appt)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Initialize appointment data
	appt.ID = uuid.New()        // Generate a new appointment ID
	appt.Status = "Unscheduled" // Default status of the appointment

	// Calculate total cost by summing the cost of all selected services
	var totalCost float64
	for _, service := range appt.Services {
		totalCost += float64(service.Cost)
	}
	appt.TotalCost = totalCost

	// Serialize services into JSON for database storage
	servicesJSON, err := json.Marshal(appt.Services)
	if err != nil {
		http.Error(w, "Failed to process services data", http.StatusInternalServerError)
		return
	}

	// Convert the appointment time (without seconds)
	// appointmentTime, err := utils.ParseISODateTimeNoSeconds(appt.DateTime)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// Insert the appointment into the database
	query := `
        INSERT INTO appointments (appointment_id, user_id, services, date, status, total_cost, created_on)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.Exec(query, appt.ID, userId, servicesJSON, appt.DateTime, appt.Status, appt.TotalCost, appt.CreatedOn)
	if err != nil {
		http.Error(w, "Failed to create appointment", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Appointment Created Successfully",
	})
}

// Get all appointments with user details
func GetAppointments(w http.ResponseWriter, r *http.Request) {
	// Define the SQL query to fetch appointments
	query := `
	 SELECT a.appointment_id, u.full_name, u.phone_no, a.services, a.total_cost, a.date, a.status
	 FROM appointments a
	 JOIN users u ON a.user_id = u.user_id
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
		var id, fullName, phoneNo, services, status, date string
		var totalCost float32

		// Scan the row into variables
		if err := rows.Scan(&id, &fullName, &phoneNo, &services, &totalCost, &date, &status); err != nil {
			log.Printf("Error scanning row: %v", err) // Log the error for debugging
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Append data to the slice
		appointments = append(appointments, map[string]interface{}{
			"id":               id,
			"full_name":        fullName,
			"phone_no":         phoneNo,
			"appointment_time": date, // Date is now a string
			"services":         services,
			"total_cost":       totalCost,
			"status":           status,
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

// Ensure the Appointment struct matches the JSON structure
type AppointmentPayload struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}

func UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	// Extract appointment ID from URL path
	if r.Method != "PATCH" {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	appointmentID := utils.GetUserID(w, r) // or use gorilla mux if needed

	// Decode the incoming request
	var appt AppointmentPayload
	err := json.NewDecoder(r.Body).Decode(&appt)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Update the appointment in the database
	query := `UPDATE appointments SET date = $1, status = $2 WHERE appointment_id = $3`
	_, err = db.Exec(query, appt.Date, appt.Status, appointmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"massage": "Appointment updated successfully"})
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
func AppointmentByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	appointmentID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var res AppointmentPayload
	query := `(
	SELECT date,status FROM appointments WHERE appointment_id = $1
	)`
	err = db.QueryRow(query, appointmentID).Scan(&res.Date, &res.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"date": res.Date, "status": res.Status,
	})
}

func GetAppointmentsByUserID(w http.ResponseWriter, r *http.Request) {
	// Get the user_id from the cookies
	userID, err := utils.GetUserIDFromCookies(r)
	if err != nil {
		http.Error(w, "User ID not found in cookies", http.StatusUnauthorized)
		return
	}

	type Appointment struct {
		ServiceNames []string `json:"service_names"`
		TotalCost    float64  `json:"total_cost"`
		CreatedOn    string   `json:"created_on"`
		Date         string   `json:"date"`
		Status       string   `json:"status"`
	}
	// SQL query to get services and other appointment details
	query := `
	SELECT 
	    jsonb_agg(s -> 'name') AS service_names,
	    total_cost,
	    created_on,
	    date,
	    status
	FROM 
	    appointments, 
	    jsonb_array_elements(appointments.services) AS s
	WHERE 
	    user_id = $1
	GROUP BY 
	    appointment_id;
	`

	// Execute the query
	rows, err := db.Query(query, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Collect the appointments
	var appointments []Appointment
	for rows.Next() {
		var appointment Appointment
		var serviceNames []byte // Store the JSONB array result

		err := rows.Scan(&serviceNames, &appointment.TotalCost, &appointment.CreatedOn, &appointment.Date, &appointment.Status)
		if err != nil {
			http.Error(w, "Error processing appointments", http.StatusInternalServerError)
			return
		}

		// Parse the JSONB array of service names
		err = json.Unmarshal(serviceNames, &appointment.ServiceNames)
		if err != nil {
			http.Error(w, "Error parsing service names", http.StatusInternalServerError)
			return
		}

		appointments = append(appointments, appointment)

		if len(appointments) == 0 {
			json.NewEncoder(w).Encode(map[string]string{"Message": "No Appointment available"})
			return
		}
	}

	// Return the appointments as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
}
