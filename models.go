package main

import (
	"database/sql"
	"log"
	"strings"
	"time"
)

// Config struct holds the database configuration details
type Config struct {
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	ServerPort string `json:"server_port"`
}

// Format struct holds the format details
type Format struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Media struct holds the media details
type Media struct {
	ID            int      `json:"id"`
	Title         string   `json:"title"`
	ArtistID      int      `json:"artist_id"`
	ArtistName    string   `json:"artist"` // This field is not stored in the database, temp field for holding the artist name to check if exists
	Media         string   `json:"media"`
	FormatID      int      `json:"format_id"`
	FormatName    string   `json:"format"` // This field is not stored in the database, temp field for holding the format name to check if exists
	DatePublished string   `json:"date_published"`
	ImageURL      string   `json:"image_url,omitempty"`
	GenreTags     []string `json:"genre_tags,omitempty"`
}

// Artist struct holds the artist details
type Artist struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	BandIDs []int  `json:"band_ids"`
}

// Band struct holds the band details
type Band struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	FormedDate time.Time `json:"formed_date"`
	Disbanded  bool      `json:"disbanded"`
	Members    []Member  `json:"members"`
}

// Member struct holds the member details
type Member struct {
	ID         int        `json:"id"`
	JoinedDate time.Time  `json:"joined_date"`
	LeftDate   *time.Time `json:"left_date,omitempty"`
}

// User struct holds the user details
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// UserMedia struct holds the details of the media owned by a user
type UserMedia struct {
	UserID   int `json:"user_id"`
	MediaID  int `json:"media_id"`
	FormatID int `json:"format_id"`
	Quantity int `json:"quantity"`
}

// NormalizeGenre normalizes the genre name based on the genre_mappings table
func NormalizeGenre(db *sql.DB, genre string) string {
	var normalizedGenre string
	err := db.QueryRow(`SELECT normalized_genre FROM genre_mappings WHERE genre = ?`, strings.ToLower(genre)).Scan(&normalizedGenre)
	if err != nil {
		if err == sql.ErrNoRows {
			return genre
		}
		log.Fatal("Failed to query genre mapping:", err)
	}
	return normalizedGenre
}
