package router

import (
	"car_service/handler"
	"net/http"
)

// func Routes() *mux.Router {
// 	r := mux.NewRouter()
// 	fileServer := http.FileServer(http.Dir("./frontend"))
// 	r.Handle("/", fileServer)
// 	// http.HandleFunc("/login", handler.Login)
// 	r.HandleFunc("/user/register", handler.RegisterUser).Methods("POST")

// 	return r
// }

func AllRoutes() {
	fileServer := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fileServer)
	http.HandleFunc("/user/register", handler.RegisterUser)
}
