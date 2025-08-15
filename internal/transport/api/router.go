package api

import (
	"bulletin-board/internal/ad"
	"github.com/gorilla/mux"
)

func NewRouter(repo ad.Repository) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/ads", GetAll(repo)).Methods("GET")
	r.HandleFunc("/ads/{id}", GetByID(repo)).Methods("GET")
	r.HandleFunc("/ads", Create(repo)).Methods("POST")
	r.HandleFunc("/ads/{id}", Update(repo)).Methods("PUT")
	r.HandleFunc("/ads/{id}", Delete(repo)).Methods("DELETE")

	return r
}
