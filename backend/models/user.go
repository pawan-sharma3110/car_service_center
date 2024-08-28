package models

import (
	"car_service/utils"
	"database/sql"
	"errors"
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

// insertUser inserts a new user into the database

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

// Validae user provided email and password
func (u User) ValidateUser(email string, password string, db *sql.DB) error {
	fmt.Println(email)
	var hashedPassword string
	query := `SELECT password FROM users WHERE email = $1`
	err:= db.QueryRow(query, u.Email).Scan(&hashedPassword)
	if err != nil {
		return fmt.Errorf("User not register with provided email: %w", err)
	}
	fmt.Println(u.Email)
	if existis := utils.UnhashPassword(password, hashedPassword); !existis {
		return errors.New("password Not match")
	}
	return nil
}
