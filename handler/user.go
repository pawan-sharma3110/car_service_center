package handler

import (
	"car_service/controller"
	"car_service/database"
	"car_service/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.DbIn()
	if err != nil {
		http.Error(w, "unable to connect database", http.StatusBadRequest)
		return
	}
	defer db.Close()
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form values
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	fullName := r.Form.Get("full_name")
	email := r.Form.Get("email")
	phoneNo := r.Form.Get("phone_no")
	password := r.Form.Get("password")
	role := r.Form.Get("role")

	id := uuid.New()
	// Create a new user
	newUser := models.User{
		ID:        id,
		FullName:  fullName,
		Email:     email,
		PhoneNo:   phoneNo,
		Password:  password,
		Role:      role,
		CreatedAt: time.Now(),
	}

	// Save the user in the database
	err = controller.InsertUser(db, w, newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]uuid.UUID{"userId": id})
}
