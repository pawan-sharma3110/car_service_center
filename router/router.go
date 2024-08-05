package router

import (
	"car_service/handler"
	"net/http"
)

func AllRoutes() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/user/register", handler.RegisterHandler)
}
