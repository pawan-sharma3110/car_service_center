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
	http.HandleFunc("/users", handler.GetAllUser)
	http.HandleFunc("/user/delete/{id}", handler.DeleteUserById)
	http.HandleFunc("/service/create", handler.CreateService)
	http.HandleFunc("/service/search", handler.SearchServicesHandler)
	http.HandleFunc("/service/get/all", handler.GetAllService)
	http.HandleFunc("/service/get", handler.GetServiceByIDHandler)
	http.HandleFunc("/service/delete/{id}", handler.DeleteById)
	http.HandleFunc("/service/delete/all", handler.AllServiceDelete)
	http.HandleFunc("/service/update/{id}", handler.UpdadeService)
	http.HandleFunc("/appointment", handler.CreateAppointment)
	http.HandleFunc("/appointment/{id}", handler.AppointmentByID)
	http.HandleFunc("/appointments", handler.GetAppointments)
	http.HandleFunc("/admin/appointment/{id}", handler.UpdateAppointment)
}
