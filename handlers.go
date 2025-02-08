package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// createRecord handles the creation of a new record
func createRecord(w http.ResponseWriter, r *http.Request) {
	var record Album
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert genre tags to a comma-separated string
	genreTags := strings.Join(record.GenreTags, ",")

	// Insert artist if it doesn't exist and get artist_id
	var artistID int
	err = db.QueryRow(`SELECT id FROM artists WHERE id = ?`, record.ArtistID).Scan(&artistID)
	if err == sql.ErrNoRows {
		http.Error(w, "Artist not found", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert record into media table
	_, err = db.Exec(`INSERT INTO media (title, date_published, image_url, genre_tags, artist_id, format_id) VALUES (?, ?, ?, ?, ?, ?)`,
		record.Title, record.DatePublished, record.ImageURL, genreTags, artistID, record.FormatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// getRecord handles retrieving a record by ID
func getRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid record ID", http.StatusBadRequest)
		return
	}

	var record Album
	var genreTags string
	err = db.QueryRow(`SELECT id, title, date_published, image_url, genre_tags, artist_id, format_id FROM media WHERE id = ?`, id).
		Scan(&record.ID, &record.Title, &record.DatePublished, &record.ImageURL, &genreTags, &record.ArtistID, &record.FormatID)
	if err != nil {
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}

	// Split genre tags string into a slice
	record.GenreTags = strings.Split(genreTags, ",")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

// updateRecord handles updating an existing record by ID
func updateRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid record ID", http.StatusBadRequest)
		return
	}

	var record Album
	err = json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert genre tags to a comma-separated string
	genreTags := strings.Join(record.GenreTags, ",")

	// Ensure the artist exists
	var artistID int
	err = db.QueryRow(`SELECT id FROM artists WHERE id = ?`, record.ArtistID).Scan(&artistID)
	if err == sql.ErrNoRows {
		http.Error(w, "Artist not found", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update the record in the media table
	_, err = db.Exec(`UPDATE media SET title = ?, date_published = ?, image_url = ?, genre_tags = ?, artist_id = ?, format_id = ? WHERE id = ?`,
		record.Title, record.DatePublished, record.ImageURL, genreTags, artistID, record.FormatID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// deleteRecord handles deleting a record by ID
func deleteRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid record ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`DELETE FROM media WHERE id = ?`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
