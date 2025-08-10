package main

import (
	rest "bulletin-board/http"
	"bulletin-board/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	store := storage.FileStore{}
	store.NewBasePath("data.json")
	r.HandleFunc("/ads", rest.AllAds(&store)).Methods("GET")
	//r.HandleFunc("/ads/{id}", rest.GetById(&store)).Methods("GET")
	r.HandleFunc("/ads", rest.Create(&store)).Methods("POST")
	//r.HandleFunc("/ads", rest.Update(&store)).Methods("PUT")
	//r.HandleFunc("/ads/{id}", rest.Delete(&store)).Methods("DELETE")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
