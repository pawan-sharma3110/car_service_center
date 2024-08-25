package controller

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// hashPassword hashes the user's password

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// insertUser inserts a new user into the database

// checkExists checks if the email or phone number already exists in the database
func CheckExists(db *sql.DB, email string, phoneNo string) error {
	var id uuid.UUID

	// Check email existence
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking email existence: %w", err)
	}
	if err == nil {
		return fmt.Errorf("provided email already registered with user ID: %v", id)
	}

	// Check phone number existence
	err = db.QueryRow("SELECT id FROM users WHERE phone_no = $1", phoneNo).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking phone number existence: %w", err)
	}
	if err == nil {
		return fmt.Errorf("provided phone number already registered with user ID: %v", id)
	}

	return nil
}
