package models

import (
	"car_service/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"time"

	"github.com/google/uuid"
)

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
}

type User struct {
	ID             uuid.UUID `json:"id"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	PhoneNo        string    `json:"phone_no"`
	Password       string    `json:"password"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	ProfilePicture []byte   `json:"profile_picture"` // New field for profile picture
	Address        Address   `json:"address"`         // New custom Address type
}
type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	ProfilePic string    `json:"profile_"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	PhoneNo    string    `json:"phone_no"`
	Role       string    `json:"role"`
	Address    Address   `json:"address"`
	CreatedAt  time.Time `json:"created_at"`
}

// insertUser inserts a new user into the database

func (u User) InsertUser(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO users (user_id, full_name, email, phone_no, password, role, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		u.ID, u.FullName, u.Email, u.PhoneNo, u.Password, u.Role, u.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert user into database: %w", err)
	}
	return nil
}

// Validae user provided email and password
func (u User) ValidateUser(email string, password string, db *sql.DB) (id uuid.UUID, role *string, err error) {

	var hashedPassword string
	query := `SELECT password,user_id,role FROM users WHERE email = $1`
	err = db.QueryRow(query, u.Email).Scan(&hashedPassword, &u.ID, &role)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf(err.Error(), u.Email)
	}
	fmt.Println(u.Email)
	fmt.Println(u.ID)

	if existis := utils.UnhashPassword(password, hashedPassword); !existis {
		return uuid.Nil, nil, errors.New("password Not match")
	}
	return u.ID, role, nil
}
func AllUsers(db *sql.DB) (users []UserResponse, err error) {
	query := `SELECT user_id, profile_picture, full_name, email, phone_no, role, address, created_at FROM users`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user UserResponse
		var profilePic sql.NullString // Handle NULL for profile picture
		var address sql.NullString    // Handle NULL for address

		// Scan the values
		err := rows.Scan(&user.ID, &profilePic, &user.FullName, &user.Email, &user.PhoneNo, &user.Role, &address, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Generate the profile picture URL
		if profilePic.Valid {
			user.ProfilePic = fmt.Sprintf("/profile-picture?user_id=%s", user.ID)
		} else {
			user.ProfilePic = "" // Handle case where there's no profile picture
		}

		// Handle address
		if address.Valid {
			var parsedAddress Address
			err := json.Unmarshal([]byte(address.String), &parsedAddress) // Assuming JSON format for address
			if err != nil {
				return nil, err
			}
			user.Address = parsedAddress
		} else {
			user.Address = Address{} // Default empty address
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}


func DeleteUser(id uuid.UUID, db *sql.DB) error {
	query := `DELETE FROM users WHERE user_id=$1`
	row, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("User not found or maybe already deleted")
	}
	return nil
}
func UpdateProfile(db *sql.DB, fullName, email, phoneNo string, address []byte, profilePicture []byte, userId uuid.UUID) error {
	query := `
		UPDATE users
		SET full_name = $1, email = $2, phone_no = $3, address = $4, profile_picture = COALESCE($5, profile_picture)
		WHERE user_id = $6
	`

	_, err := db.Exec(query, fullName, email, phoneNo, address, profilePicture, userId)
	if err != nil {
		return err
	}

	return nil
}
