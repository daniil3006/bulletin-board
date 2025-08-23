package api

import (
	"bulletin-board/internal/middleware"
	"github.com/gorilla/mux"
)

func (h Handler) NewRouter(r *mux.Router) {
	r.HandleFunc("/ads", h.GetAll()).Methods("GET")
	r.HandleFunc("/ads/{id}", h.GetByID()).Methods("GET")

	secured := r.PathPrefix("/ads").Subrouter()
	secured.Use(middleware.AuthMiddleware("iuNvi8sa5oiHOajKfn93hFL93gb"))

	secured.HandleFunc("", h.Create()).Methods("POST")
	secured.HandleFunc("/{id}", h.Update()).Methods("PUT")
	secured.HandleFunc("/{id}", h.Delete()).Methods("DELETE")
}
