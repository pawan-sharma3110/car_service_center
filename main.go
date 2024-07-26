package main

import (
	"car_service/db"
	"car_service/router"
	"fmt"

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
	// r := router.Routes()
	router.AllRoutes()

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
