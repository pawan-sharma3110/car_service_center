package main

import (
	"car_service/router"
	"fmt"
	"os"

	"log"
	"net/http"
)

func main() {

	// config.Load()
	// Check if environment variables are loaded
	log.Printf("DB_USER: %s", os.Getenv("DB_USER"))
	log.Printf("DB_PASSWORD: %s", os.Getenv("DB_PASSWORD"))
	log.Printf("DB_HOST: %s", os.Getenv("DB_HOST"))
	log.Printf("DB_PORT: %s", os.Getenv("DB_PORT"))
	log.Printf("DB_NAME: %s", os.Getenv("DB_NAME"))
	log.Printf("SSL_MODE: %s", os.Getenv("SSL_MODE"))
	router.AllRoutes()

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
