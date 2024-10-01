package main

import (
	"car_service/router"
	"fmt"

	"log"
	"net/http"
)

func main() {

	// config.Load()
	router.AllRoutes()

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
