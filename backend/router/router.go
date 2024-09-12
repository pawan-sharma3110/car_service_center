package router

import (
	"car_service/handler"
	"car_service/middleware"
	"net/http"
)

func AllRoutes() {
	fileServer := http.FileServer(http.Dir("../static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/user/register", handler.RegisterHandler)
	http.HandleFunc("/user/login", handler.Login)
	http.HandleFunc("/service/search", handler.SearchServicesHandler)
	http.HandleFunc("/service/get/all", handler.GetAllService)
	http.HandleFunc("/appointment", handler.CreateAppointment)
	
	
	// Admin Services


	http.HandleFunc("/users", middleware.RoleBasedAuth(handler.GetAllUser, "admin"))
	http.HandleFunc("/user/delete/{id}", middleware.RoleBasedAuth(handler.DeleteUserById, "admin"))
	http.HandleFunc("/service/create", middleware.RoleBasedAuth(handler.CreateService, "admin"))
	http.HandleFunc("/service/get", middleware.RoleBasedAuth(handler.GetServiceByIDHandler, "admin"))
	http.HandleFunc("/service/delete/{id}", middleware.RoleBasedAuth(handler.DeleteById, "admin"))
	http.HandleFunc("/service/delete/all", middleware.RoleBasedAuth(handler.AllServiceDelete, "admin"))
	http.HandleFunc("/service/update/{id}", middleware.RoleBasedAuth(handler.UpdadeService, "admin"))
	http.HandleFunc("/appointment/{id}", middleware.RoleBasedAuth(handler.AppointmentByID, "admin"))
	http.HandleFunc("/appointments", middleware.RoleBasedAuth(handler.GetAppointments, "admin"))
	http.HandleFunc("/admin/appointment/{id}", middleware.RoleBasedAuth(handler.UpdateAppointment,"admin"))
}
