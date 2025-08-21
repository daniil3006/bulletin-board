package api

import "github.com/gorilla/mux"

func (h Handler) NewRouter(r *mux.Router) {
	r.HandleFunc("/users", h.GetAll()).Methods("GET")
	r.HandleFunc("/users/{id}", h.GetByID()).Methods("GET")
	r.HandleFunc("/users/{id}/ads", h.GetUsersAds()).Methods("GET")
	r.HandleFunc("/users", h.Create()).Methods("POST")
	r.HandleFunc("/users/{id}", h.Update()).Methods("PUT")
	r.HandleFunc("/users/{id}", h.Delete()).Methods("DELETE")
}
