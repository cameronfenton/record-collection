package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	initDB()
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/records", createRecord).Methods("POST")
	router.HandleFunc("/records/{id}", getRecord).Methods("GET")
	router.HandleFunc("/records/{id}", updateRecord).Methods("PUT")
	router.HandleFunc("/records/{id}", deleteRecord).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
