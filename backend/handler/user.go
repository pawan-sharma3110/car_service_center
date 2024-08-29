package handler

import (
	"car_service/controller"
	"car_service/database"
	"car_service/models"
	"car_service/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var db *sql.DB = database.DbIn()

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	// Ensure the request is a POST request
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if user.FullName == "" || user.Email == "" || user.PhoneNo == "" || user.Password == "" || user.Role == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Check if email or phone number already exists
	err = controller.CheckExists(db, user.Email, user.PhoneNo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a new UUID for the user
	id := uuid.New()

	// Create a new user struct
	newUser := models.User{
		ID:        id,
		FullName:  user.FullName,
		Email:     user.Email,
		PhoneNo:   user.PhoneNo,
		Password:  user.Password,
		Role:      user.Role,
		CreatedAt: time.Now(),
	}

	// Hash the password
	newUser.Password, err = utils.HashPassword(newUser.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Insert the new user into the database
	err = newUser.InsertUser(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Respond with success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]uuid.UUID{"userId": newUser.ID})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	_, err = user.ValidateUser(user.Email, user.Password, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := utils.GernateJwt(user.Email, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
