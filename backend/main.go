package main

import (
	"car_service/backend/db"
	"car_service/backend/router"
	"log"
	"net/http"
)

func main() {
	database, err := db.DbIn()
	if err != nil {
		panic(err)
	}
	defer database.Close()

	// Initialize the router
	r := router.Routes()

	// Start the server
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
