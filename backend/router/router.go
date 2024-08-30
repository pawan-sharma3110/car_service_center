package router

import (
	"car_service/handler"
	"net/http"
)

func AllRoutes() {
	fileServer := http.FileServer(http.Dir("../static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/user/register", handler.RegisterHandler)
	http.HandleFunc("/user/login", handler.Login)

	// Admin Services

	http.HandleFunc("/service/create", handler.CreateService)
	http.HandleFunc("/service/get/all", handler.GetAllService)
	http.HandleFunc("/service/deletebyid/{id}", handler.DeleteById)
	http.HandleFunc("/service/delete/all", handler.AllServiceDelete)
}
