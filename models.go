package main

import (
	"time"
)

type Record struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	DatePublished time.Time `json:"date_published"`
	ImageURL      string    `json:"image_url"`
	GenreTags     []string  `json:"genre_tags"`
	ArtistID      int       `json:"artist_id"`
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
