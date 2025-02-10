package main

import (
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

// Album struct holds the format details
type Format struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Media struct holds the media details
type Media struct {
	ID            int      `json:"id"`
	Title         string   `json:"title"`
	ArtistID      string   `json:"artist_id"`
	Artist        string   `json:"artist"` // This field is not stored in the database, temp field for holding the artist name to check if exists
	Media         string   `json:"media"`
	FormatID      int      `json:"format_id"`
	Format        string   `json:"format"` // This field is not stored in the database, temp field for holding the format name to check if exists
	DatePublished string   `json:"date_published"`
	ImageURL      string   `json:"image_url,omitempty"`
	GenreTags     []string `json:"genre_tags,omitempty"`
}

type Artist struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	BandIDs []int  `json:"band_ids"`
}

type Band struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	FormedDate time.Time `json:"formed_date"`
	Disbanded  bool      `json:"disbanded"`
	Members    []Member  `json:"members"`
}

type Member struct {
	ID         int        `json:"id"`
	JoinedDate time.Time  `json:"joined_date"`
	LeftDate   *time.Time `json:"left_date,omitempty"`
}
