package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// createMedia handles the creation of a new media
func createMedia(w http.ResponseWriter, r *http.Request) {
	var m Media
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert genre tags to a comma-separated string
	genreTags := strings.Join(m.GenreTags, ",")

	// Insert artist if it doesn't exist and get artist_id
	var artistID int
	err = db.QueryRow(`SELECT id FROM artists WHERE id = ?`, m.ArtistID).Scan(&artistID)
	if err == sql.ErrNoRows {
		http.Error(w, "Artist not found", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert media into media table
	_, err = db.Exec(`INSERT INTO media (title, date_published, image_url, genre_tags, artist_id, format_id) VALUES (?, ?, ?, ?, ?, ?)`,
		m.Title, m.DatePublished, m.ImageURL, genreTags, artistID, m.FormatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// getMedia handles retrieving all media
func getMedia(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
        SELECT 
            m.id, m.title, m.date_published, m.image_url, m.genre_tags, 
            m.artist_id, a.name, m.format_id, f.name
        FROM media m 
        JOIN artists a ON m.artist_id = a.id
        JOIN formats f ON m.format_id = f.id
    `)
	if err != nil {
		http.Error(w, "Failed to retrieve media", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var media []Media

	for rows.Next() {
		var m Media
		var genreTags string
		if err := rows.Scan(
			&m.ID, &m.Title, &m.DatePublished, &m.ImageURL, &genreTags,
			&m.ArtistID, &m.ArtistName, &m.FormatID, &m.FormatName,
		); err != nil {
			http.Error(w, "Failed to scan media", http.StatusInternalServerError)
			return
		}
		// Split genre tags string into a slice
		m.GenreTags = strings.Split(genreTags, ",")
		media = append(media, m)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating over media", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(media)
}

// getMediaById handles retrieving a media by ID
func getMediaById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	var m Media
	var genreTags string
	err = db.QueryRow(`
        SELECT 
            m.id, m.title, m.date_published, m.image_url, m.genre_tags, 
            m.artist_id, a.name, m.format_id, f.name
        FROM media m 
        JOIN artists a ON m.artist_id = a.id
        JOIN formats f ON m.format_id = f.id
        WHERE m.id = ?
    `, id).Scan(
		&m.ID, &m.Title, &m.DatePublished, &m.ImageURL, &genreTags,
		&m.ArtistID, &m.ArtistName, &m.FormatID, &m.FormatName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Media not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve media", http.StatusInternalServerError)
		}
		return
	}

	// Split genre tags string into a slice
	m.GenreTags = strings.Split(genreTags, ",")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

// updateMedia handles updating an existing media by ID
func updateMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	var m Media
	err = json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert genre tags to a comma-separated string
	genreTags := strings.Join(m.GenreTags, ",")

	// Ensure the artist exists
	var artistID int
	err = db.QueryRow(`SELECT id FROM artists WHERE id = ?`, m.ArtistID).Scan(&artistID)
	if err == sql.ErrNoRows {
		http.Error(w, "Artist not found", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update the media in the media table
	_, err = db.Exec(`UPDATE media SET title = ?, date_published = ?, image_url = ?, genre_tags = ?, artist_id = ?, format_id = ? WHERE id = ?`,
		m.Title, m.DatePublished, m.ImageURL, genreTags, artistID, m.FormatID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// deleteMedia handles deleting a media by ID
func deleteMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`DELETE FROM media WHERE id = ?`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
