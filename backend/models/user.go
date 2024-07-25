package models

import "time"

type User struct {
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	PhoneNo   string    `json:"phone_no"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
