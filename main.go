package main

import (
	"car_service/handler"
	"fmt"

	"log"
	"net/http"
)

func main() {
	// http.FileServer(http.Dir("./forntend"))
	fileServer := http.FileServer(http.Dir("../static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/user/register", handler.RegisterHandler)
	// router.AllRoutes()

	fmt.Println("Starting server at port 8001")
	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatal(err)
	}
}
