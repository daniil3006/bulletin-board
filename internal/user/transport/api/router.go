package api

import (
	"bulletin-board/internal/middleware"
	"github.com/gorilla/mux"
)

func (h Handler) NewRouter(r *mux.Router) {
	r.HandleFunc("/users", h.Create()).Methods("POST")
	r.HandleFunc("/sign-in", h.SignIn()).Methods("POST")
	r.HandleFunc("/users", h.GetAll()).Methods("GET")
	r.HandleFunc("/users/{id}", h.GetByID()).Methods("GET")
	r.HandleFunc("/users/{id}/ads", h.GetUsersAds()).Methods("GET")

	secured := r.PathPrefix("/users").Subrouter()
	secured.Use(middleware.AuthMiddleware("iuNvi8sa5oiHOajKfn93hFL93gb"))

	secured.HandleFunc("/{id}", h.Update()).Methods("PUT")
	secured.HandleFunc("/{id}", h.Delete()).Methods("DELETE")
}
