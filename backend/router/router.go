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

	http.HandleFunc("/profile-picture", middleware.RoleBasedAuth(handler.GetProfilePicture, "admin"))
	http.HandleFunc("/user/profile/{id}", middleware.RoleBasedAuth(handler.GetUserProfile, "admin", "user"))
	http.HandleFunc("/user/update-profile/{id}", middleware.RoleBasedAuth(handler.UpdateUserProfile, "admin", "user"))
	http.HandleFunc("/user/logout", handler.LogoutHandler)
	http.HandleFunc("/service/search", handler.SearchServicesHandler)
	http.HandleFunc("/service/get/all", middleware.RoleBasedAuth(handler.GetAllService, "admin", "user"))
	http.HandleFunc("/appointment", middleware.RoleBasedAuth(handler.CreateAppointment, "admin", "user"))
	http.HandleFunc("/appointment/by_id", middleware.RoleBasedAuth(handler.GetAppointmentsByUserID, "admin", "user"))
	// Admin Services

	http.HandleFunc("/users", middleware.RoleBasedAuth(handler.GetAllUser, "admin"))
	http.HandleFunc("/profile-picture/{id}", middleware.RoleBasedAuth(handler.GetProfilePicture, "admin"))
	http.HandleFunc("/user/delete/{id}", middleware.RoleBasedAuth(handler.DeleteUserById, "admin"))
	http.HandleFunc("/service/create", middleware.RoleBasedAuth(handler.CreateService, "admin"))
	http.HandleFunc("/service/get", middleware.RoleBasedAuth(handler.GetServiceByIDHandler, "admin"))
	http.HandleFunc("/service/delete/{id}", middleware.RoleBasedAuth(handler.DeleteById, "admin"))
	http.HandleFunc("/service/delete/all", middleware.RoleBasedAuth(handler.AllServiceDelete, "admin"))
	http.HandleFunc("/service/update/{id}", middleware.RoleBasedAuth(handler.UpdadeService, "admin"))
	http.HandleFunc("/appointment/{id}", middleware.RoleBasedAuth(handler.AppointmentByID, "admin"))
	http.HandleFunc("/appointments", middleware.RoleBasedAuth(handler.GetAppointments, "admin"))
	http.HandleFunc("/admin/appointment/{id}", middleware.RoleBasedAuth(handler.UpdateAppointment, "admin"))
}
