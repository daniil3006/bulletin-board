package api

import (
	"github.com/gorilla/mux"
)

func (h Handler) NewRouter(r *mux.Router) {
	r.HandleFunc("/ads", h.GetAll()).Methods("GET")
	r.HandleFunc("/ads/{id}", h.GetByID()).Methods("GET")
	r.HandleFunc("/ads", h.Create()).Methods("POST")
	r.HandleFunc("/ads/{id}", h.Update()).Methods("PUT")
	r.HandleFunc("/ads/{id}", h.Delete()).Methods("DELETE")
}
