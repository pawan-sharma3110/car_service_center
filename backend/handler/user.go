package handler

import (
	"car_service/controller"
	"car_service/database"
	"car_service/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var db *sql.DB = database.DbIn()

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	defer db.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form values
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	fullName := r.Form.Get("full_name")
	email := r.Form.Get("email")
	phoneNo := r.Form.Get("phone_no")
	password := r.Form.Get("password")
	role := r.Form.Get("role")

	fmt.Println("Form values:")
	fmt.Println("Full Name:", r.Form.Get("full_name"))
	fmt.Println("Email:", r.Form.Get("email"))
	fmt.Println("Phone No:", r.Form.Get("phone_no"))
	fmt.Println("Password:", r.Form.Get("password"))
	fmt.Println("Role:", r.Form.Get("role"))

	if fullName == "" || email == "" || phoneNo == "" || password == "" || role == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	// Check if email or phone number already exists
	err = controller.CheckExists(db, email, phoneNo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Generate a new UUID for the user
	id := uuid.New()

	// Create a new user struct
	newUser := models.User{
		ID:        id,
		FullName:  fullName,
		Email:     email,
		PhoneNo:   phoneNo,
		Password:  password,
		Role:      role,
		CreatedAt: time.Now(),
	}

	// Hash the password
	newUser.Password, err = controller.HashPassword(newUser.Password)
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

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]uuid.UUID{"userId": newUser.ID})
}
