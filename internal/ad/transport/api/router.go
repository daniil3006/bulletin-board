package api

import (
	"bulletin-board/internal/middleware"
	"github.com/gorilla/mux"
	"os"
)

func (h Handler) NewRouter(r *mux.Router) {
	r.HandleFunc("/ads", h.GetAll()).Methods("GET")
	r.HandleFunc("/ads/{id}", h.GetByID()).Methods("GET")

	secretKey := os.Getenv("SINGING_KEY")
	secured := r.PathPrefix("/ads").Subrouter()
	secured.Use(middleware.AuthMiddleware(secretKey))

	secured.HandleFunc("", h.Create()).Methods("POST")
	secured.HandleFunc("/{id}", h.Update()).Methods("PUT")
	secured.HandleFunc("/{id}", h.Delete()).Methods("DELETE")
}
