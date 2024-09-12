package utils

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = "self_project"

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

func GernateJwt(email string, id uuid.UUID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
	})
	return token.SignedString([]byte(jwtKey))
}

func RoleBasedAuth(next http.HandlerFunc, allowedRoles ...string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userRole := (*claims)["role"].(string)

		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func FormatDate(date sql.NullTime) interface{} {
	if date.Valid {
		return date.Time.Format(time.RFC3339) // Format date as RFC3339 string
	}
	return nil
}
