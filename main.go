package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	initDB()
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/media", createMedia).Methods("POST")
	router.HandleFunc("/media", getMedia).Methods("GET")
	router.HandleFunc("/media/{id}", getMediaById).Methods("GET")
	router.HandleFunc("/media/{id}", updateMedia).Methods("PUT")
	router.HandleFunc("/media/{id}", deleteMedia).Methods("DELETE")

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Add the CORS middleware to the router
	handler := c.Handler(router)

	log.Printf("Starting server on port %s...", config.ServerPort)
	log.Fatal(http.ListenAndServe(":"+config.ServerPort, handler)) // Use handler instead of router
}
