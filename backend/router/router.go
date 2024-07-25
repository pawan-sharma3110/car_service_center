package router

import (
	"car_service/backend/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	r := mux.NewRouter()
	fileServer := http.FileServer(http.Dir("../frontend"))
	r.Handle("/", fileServer)
	r.HandleFunc("/user/register", handler.RegisterUser).Methods("POST")

	return r
}

// func Routes() *http.HandlerFunc {
// 	r := http.HandlerFunc()
// }
