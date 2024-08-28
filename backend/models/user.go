package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	PhoneNo   string    `json:"phone_no"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}


func (u User) InsertUser(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO users (id, full_name, email, phone_no, password, role, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		u.ID, u.FullName, u.Email, u.PhoneNo, u.Password, u.Role, u.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert user into database: %w", err)
	}
	return nil
}
