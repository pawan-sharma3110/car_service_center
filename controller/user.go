package controller

import (
	"car_service/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
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
		return fmt.Errorf("error creating users table: %w", err)
	}
	err = checkExists(db, user.Email, user.PhoneNo)
	if err != nil {
		return err
	}
	// Hash the password

	user.Password, err = converHash(user.Password)
	if err != nil {
		return err
	}

	// Save the user to the database
	_, err = db.Exec(`
		INSERT INTO users (id, full_name, email, phone_no, password, role, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.ID, user.FullName, user.Email, user.PhoneNo, user.Password, user.Role, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("falid to insert user in database : %w", err)
	}

	return nil
}

func checkExists(db *sql.DB, email string, phoneNo string) error {
	// Check if the email already exists
	var id uuid.UUID
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking email existence: %w", err)
	}
	if err == nil {
		return fmt.Errorf("provided email already registered with this user ID: %v", id)
	}

	// Check if the phone number already exists
	err = db.QueryRow("SELECT id FROM users WHERE phone_no = $1", phoneNo).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking phone number existence: %w", err)
	}
	if err == nil {
		return fmt.Errorf("provided mobile number already registered with this user ID: %v", id)
	}

	return nil
}

func converHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", fmt.Errorf("failed to convert password in hash %w", err)
	}
	return string(hashedPassword), nil
}
