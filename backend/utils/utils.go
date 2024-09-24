package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var JwtKey = "self_project"

// hashPassword hashes the user's password

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// compaire user's password
func UnhashPassword(plainPassword string, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainPassword))
	return err == nil
}

//gernate jwt

func GenerateJwt(email string, id uuid.UUID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour.Abs() * 15).Unix(),
	})
	return token.SignedString([]byte(JwtKey))
}

func ParseISODateTimeNoSeconds(dateTimeStr string) (time.Time, error) {
	// Layout for ISO 8601 without seconds
	layout := "2006-01-02T15:04"

	// Parse the provided string based on the layout
	parsedTime, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		return time.Time{}, errors.New("invalid date/time format, expected ISO 8601 (YYYY-MM-DDTHH:MM)")
	}
	return parsedTime, nil
}

func GetUserID(w http.ResponseWriter, r *http.Request) uuid.UUID {
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return uuid.Nil
	}

	userId, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return uuid.Nil
	}
	return userId
}
func GetUserIDFromCookies(r *http.Request) (uuid.UUID, error) {
	// Retrieve the cookie named "user_id"
	cookie, err := r.Cookie("user_id")
	if err != nil {
		if err == http.ErrNoCookie {
			return uuid.Nil, errors.New("missing user ID cookie")
		}

		return uuid.Nil, err
	}

	// Return the user ID from the cookie value
	userId := cookie.Value
	if userId == "" {
		return uuid.Nil, errors.New("user ID not found in cookie")
	}
	id, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
