package utils

import (
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

// Helper function to format time to "YYYY-MM-DD hh:mm:ss AM/PM"
func FormatToAMPM(t time.Time) string {
	return t.Format("2006-01-02 03:04:05 PM")
}

// Helper function to parse time from UI (AM/PM) and convert to local time
func ParseFromAMPM(dateTime string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local") // Use local time from the device
	t, err := time.ParseInLocation("2006-01-02 03:04:05 PM", dateTime, loc)
	return t, err
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
