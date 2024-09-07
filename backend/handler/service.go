package handler

import (
	"car_service/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func CreateService(w http.ResponseWriter, r *http.Request) {
	// Ensure the request is a POST request
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var service models.Service
	err := json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if service.Name == "" || service.Cost == 0 || service.Description == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Set additional fields before saving
	service.ServiceID = uuid.New()
	service.CreatedAt = time.Now()

	// Save service to the database (replace with your actual db interaction)
	err = service.InsertService(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created service ID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]uuid.UUID{"service_id": service.ServiceID})
}
func GetServiceByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the service_id from the URL query parameter
	serviceIDStr := r.URL.Query().Get("id")
	if serviceIDStr == "" {
		http.Error(w, "Missing service ID", http.StatusBadRequest)
		return
	}

	// Parse the UUID from the string
	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	// Query the database for the service by its ID
	service := models.GetServiceByID(db, serviceID, w)
	// Return the service as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(service)
}

func GetAllService(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	services, err := models.AllService(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if services == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"services": "No services found in database"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services)
}
func DeleteById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	serviceID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}

	err = models.DeleteService(serviceID, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Respond to the client

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service with ID %s deleted successfully", serviceID)
}

func AllServiceDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	err := models.DeleteAll(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, " All service deleted successfully")
}
func UpdadeService(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	sID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}
	var service models.Service
	err = json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updated := models.Service{
		ServiceID:   sID,
		Name:        service.Name,
		Description: service.Description,
		Cost:        service.Cost,
	}
	err = updated.Update(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service updated successfully")
}

func SearchServicesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	searchQuery := "%" + strings.TrimSpace(query) + "%"
	rows, err := db.Query(`
        SELECT service_id, name, cost, description
        FROM services
        WHERE name ILIKE $1 OR cost::text ILIKE $1 OR service_id::text ILIKE $1
        `, searchQuery)
	if err != nil {
		http.Error(w, "Failed to query the database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// services := []struct {
	//     ServiceID   string  `json:"service_id"`
	//     Name        string  `json:"name"`
	//     Cost        float32 `json:"cost"`
	//     Description string  `json:"description"`
	// }{}
	var services []models.Service
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.ServiceID, &s.Name, &s.Cost, &s.Description); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		services = append(services, s)
	}

	if len(services) == 0 {
		http.Error(w, "No services found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}
