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

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserID(w, r) // Assuming this function retrieves user ID from the cookie

	// Query the user details from the database
	var user models.User
	var addressData []byte // Variable to hold the raw JSON data for the address

	err := db.QueryRow(`SELECT user_id, full_name, email, phone_no, role, address FROM users WHERE user_id = $1`, userID).
		Scan(&user.ID, &user.FullName, &user.Email, &user.PhoneNo, &user.Role, &addressData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Unmarshal the address JSON data into the Address struct
	if err := json.Unmarshal(addressData, &user.Address); err != nil {
		http.Error(w, "Error parsing address data", http.StatusInternalServerError)
		return
	}

	// Send the user data as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
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
func GetProfilePicture(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL parameters
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Query the database for the profile picture
	var profilePicture []byte
	err := db.QueryRow(`SELECT profile_picture FROM users WHERE user_id = $1`, userID).Scan(&profilePicture)
	if err != nil {
		http.Error(w, "Error retrieving profile picture", http.StatusInternalServerError)
		return
	}

	// Set the appropriate content type (assuming you're storing common image types)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(profilePicture) // Write the binary image data to the response
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

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID
	userID := utils.GetUserID(w, r)

	// Parse multipart form data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}

	// Read profile picture
	var profilePicture []byte
	hasProfilePicture := false
	file, _, err := r.FormFile("profile_picture")
	if err == nil {
		defer file.Close()
		profilePicture, _ = io.ReadAll(file)
		hasProfilePicture = true
	}

	// Extract other form fields
	fullName := r.FormValue("full_name")
	email := r.FormValue("email")
	phoneNo := r.FormValue("phone_no")
	role := r.FormValue("role")
	password := r.FormValue("password") // New field for password

	// Extract and parse address JSON
	var addressJSON *string
	addressStr := r.FormValue("address")
	if addressStr != "" {
		addressJSONStr := addressStr
		addressJSON = &addressJSONStr
	}

	// Prepare SQL query and arguments
	var query string
	var args []interface{}

	// Check if password field is empty or not
	if password != "" {
		// Hash the new password (assuming you have a helper for password hashing)
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		if hasProfilePicture {
			query = `
                UPDATE users
                SET full_name = COALESCE(NULLIF($2, ''), full_name),
                    email = COALESCE(NULLIF($3, ''), email),
                    phone_no = COALESCE(NULLIF($4, ''), phone_no),
                    profile_picture = $5,
                    address = COALESCE(NULLIF($6::jsonb, '{}'::jsonb), address),
                    role = COALESCE(NULLIF($7, ''), role),
                    password = $8
                WHERE user_id = $1
            `
			args = []interface{}{userID, fullName, email, phoneNo, profilePicture, addressJSON, role, hashedPassword}
		} else {
			query = `
                UPDATE users
                SET full_name = COALESCE(NULLIF($2, ''), full_name),
                    email = COALESCE(NULLIF($3, ''), email),
                    phone_no = COALESCE(NULLIF($4, ''), phone_no),
                    address = COALESCE(NULLIF($5::jsonb, '{}'::jsonb), address),
                    role = COALESCE(NULLIF($6, ''), role),
                    password = $7
                WHERE user_id = $1
            `
			args = []interface{}{userID, fullName, email, phoneNo, addressJSON, role, hashedPassword}
		}
	} else {
		// Password field is empty, retain the old password
		if hasProfilePicture {
			query = `
                UPDATE users
                SET full_name = COALESCE(NULLIF($2, ''), full_name),
                    email = COALESCE(NULLIF($3, ''), email),
                    phone_no = COALESCE(NULLIF($4, ''), phone_no),
                    profile_picture = $5,
                    address = COALESCE(NULLIF($6::jsonb, '{}'::jsonb), address),
                    role = COALESCE(NULLIF($7, ''), role)
                WHERE user_id = $1
            `
			args = []interface{}{userID, fullName, email, phoneNo, profilePicture, addressJSON, role}
		} else {
			query = `
                UPDATE users
                SET full_name = COALESCE(NULLIF($2, ''), full_name),
                    email = COALESCE(NULLIF($3, ''), email),
                    phone_no = COALESCE(NULLIF($4, ''), phone_no),
                    address = COALESCE(NULLIF($5::jsonb, '{}'::jsonb), address),
                    role = COALESCE(NULLIF($6, ''), role)
                WHERE user_id = $1
            `
			args = []interface{}{userID, fullName, email, phoneNo, addressJSON, role}
		}
	}

	// Execute the query
	_, err = db.Exec(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User profile updated successfully"))
}
