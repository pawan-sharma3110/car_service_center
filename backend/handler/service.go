package handler

import (
	"car_service/models"
	"encoding/json"
	"fmt"
	"net/http"
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
