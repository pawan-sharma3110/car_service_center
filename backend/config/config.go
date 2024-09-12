package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	SSLMode    string
	
)

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBName = os.Getenv("DB_NAME")
	SSLMode = os.Getenv("SSL_MODE")
}
