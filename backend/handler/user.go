package handler

import (
	"car_service/controller"
	"car_service/database"
	"car_service/models"
	"car_service/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
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

	id, role, err := user.ValidateUser(user.Email, user.Password, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token, err := utils.GenerateJwt(user.Email, id, *role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set the token and user_id as cookies
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
		// HttpOnly: true, // Keep HttpOnly for security
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: id.String(),
		// HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	// Send success response including the token in the body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   token, // Add token to the response
		"role":    *role,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-24 * time.Hour),
		Path:    "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   "",
		Expires: time.Now().Add(-24 * time.Hour),
		Path:    "/",
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// var users models.UserResponse
	users, err := models.AllUsers(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
}
func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userId := utils.GetUserID(w, r)
	err := models.DeleteUser(userId, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User with ID %s deleted successfully", userId)
}

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	userId := utils.GetUserID(w, r)

	// Parse the multipart form to handle file uploads
	r.ParseMultipartForm(3 << 20) // Limit file size to 3MB

	// Get user profile data from the form
	fullName := r.FormValue("full_name")
	email := r.FormValue("email")
	phoneNo := r.FormValue("phone_no")
	address := models.Address{
		Street:  r.FormValue("street"),
		City:    r.FormValue("city"),
		State:   r.FormValue("state"),
		ZipCode: r.FormValue("zip_code"),
	}

	// Convert address to JSON for storage in DB
	add, err := json.Marshal(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process the profile picture file if uploaded
	var profilePicture []byte
	file, handler, err := r.FormFile("profile_picture")
	if err == nil && handler != nil {
		defer file.Close()

		// Check the file extension
		fileExt := filepath.Ext(handler.Filename)
		if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" {
			http.Error(w, "Only JPG, JPEG, and PNG files are allowed", http.StatusBadRequest)
			return
		}

		// Read the file into a byte slice (to store it as binary in the database)
		profilePicture, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading image file", http.StatusInternalServerError)
			return
		}
	}

	// Update user profile in the database
	err = models.UpdateProfile(db, w, fullName, email, phoneNo, add, profilePicture, userId)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}
