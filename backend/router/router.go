package router

import (
	"car_service/backend/handler"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/user/register", handler.RegisterUser)
	return r
}
