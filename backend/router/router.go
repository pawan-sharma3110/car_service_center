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
	http.HandleFunc("/user/appointments/create", handler.CreateAppointment)

	// Admin Services

	http.HandleFunc("/service/create", handler.CreateService)
	http.HandleFunc("/service/search", handler.SearchServicesHandler)
	http.HandleFunc("/service/get/all", handler.GetAllService)
	http.HandleFunc("/service/get", handler.GetServiceByIDHandler)
	http.HandleFunc("/service/delete/{id}", handler.DeleteById)
	http.HandleFunc("/service/delete/all", handler.AllServiceDelete)
	http.HandleFunc("/service/update/{id}", handler.UpdadeService)
	http.HandleFunc("/appointments", handler.GetAllAppointmentsHandler)
	http.HandleFunc("/admin/appointment/{id}", handler.AppointmentStatusUpdate)
}
