package main

import (
	"car_service/router"
	"fmt"

	"log"
	"net/http"
)

func main() {

	router.AllRoutes()

	fmt.Println("Starting server at port 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
