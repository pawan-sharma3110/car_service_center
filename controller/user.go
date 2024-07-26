package controller

import (
	"car_service/models"
	"database/sql"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func InsertUser(db *sql.DB, w http.ResponseWriter, user models.User) error {
	// Create the users table if it doesn't exist
	query := `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		full_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		phone_no VARCHAR(20) UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error creating users table: %v", err)
		http.Error(w, "Error creating users table", http.StatusInternalServerError)
		return err
	}

	// Check if the user already exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1 OR phone_no = $2", user.Email, user.PhoneNo).Scan(&count)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return err
	}

	if count > 0 {
		log.Printf("Email or phone number already exists: email=%s, phone_no=%s", user.Email, user.PhoneNo)
		http.Error(w, "Email or phone number already exists", http.StatusConflict)
		return err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return err
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()

	// Save the user to the database
	_, err = db.Exec(`
		INSERT INTO users (id, full_name, email, phone_no, password, role, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.ID, user.FullName, user.Email, user.PhoneNo, user.Password, user.Role, user.CreatedAt)
	if err != nil {
		log.Printf("Error saving user: %v", err)
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return err
	}

	log.Printf("User registered successfully: id=%s, email=%s, phone_no=%s", user.ID, user.Email, user.PhoneNo)
	return nil
}
