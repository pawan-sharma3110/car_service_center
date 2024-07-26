package handler

import (
	"car_service/controller"
	"car_service/models"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

var Database *sql.DB

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/registers" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Parse form values
	newUser := models.User{
		ID:       uuid.New(),
		FullName: r.FormValue("full_name"),
		Email:    r.FormValue("email"),
		PhoneNo:  r.FormValue("phone_no"),
		Password: r.FormValue("password"),
		Role:     r.FormValue("role"),
	}

	// Validate required fields
	if newUser.FullName == "" || newUser.Email == "" || newUser.PhoneNo == "" || newUser.Password == "" || newUser.Role == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Insert the user into the database
	err := controller.InsertUser(Database, w, newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"}); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
