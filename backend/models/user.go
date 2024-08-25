package models

import (
	"database/sql"
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

func CreateUser(db *sql.DB, user User) error {
	_, err := db.Exec("INSERT INTO users (id, full_name, email, phone_no, password, role, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.ID, user.FullName, user.Email, user.PhoneNo, user.Password, user.Role, user.CreatedAt)
	return err
}
