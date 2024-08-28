package controller

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)



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
