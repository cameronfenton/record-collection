package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	initDB()
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/records", createRecord).Methods("POST")
	router.HandleFunc("/records/{id}", getRecord).Methods("GET")
	router.HandleFunc("/records/{id}", updateRecord).Methods("PUT")
	router.HandleFunc("/records/{id}", deleteRecord).Methods("DELETE")

	log.Printf("Starting server on port %s...", config.ServerPort)
	log.Fatal(http.ListenAndServe(":"+config.ServerPort, router))
}
