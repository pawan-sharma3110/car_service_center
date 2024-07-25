package handler

import (
	"car_service/backend/controller"
	"car_service/backend/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

var Database *sql.DB

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/registers" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	var newUser models.User
	newUser.ID = uuid.New()
	newUser.FullName = r.FormValue("full_name")
	newUser.Email = r.FormValue("email")
	newUser.PhoneNo = r.FormValue("phone_no")
	newUser.Password = r.FormValue("password")
	newUser.Role = r.FormValue("role")

	// Handle your form data (e.g., save to database, validate, etc.)
	// For demonstration, we'll just print them
	fmt.Fprintf(w, "ID: %s\n", newUser.ID)
	fmt.Fprintf(w, "Full Name: %s\n", newUser.FullName)
	fmt.Fprintf(w, "Email: %s\n", newUser.Email)
	fmt.Fprintf(w, "Phone Number: %s\n", newUser.PhoneNo)
	fmt.Fprintf(w, "Password: %s\n", newUser.Password)
	fmt.Fprintf(w, "Role: %s\n", newUser.Role)
	err := controller.InsertUser(Database, w, newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
