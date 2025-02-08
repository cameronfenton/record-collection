package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// db is the global database connection pool
var db *sql.DB

// loadConfig reads the configuration from a JSON file
func loadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// connectToMySQL connects to the MySQL server
func connectToMySQL(config *Config) error {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/", config.DBUser, config.DBPassword, config.DBHost, config.DBPort)
	var err error
	db, err = sql.Open("mysql", connStr)
	return err
}

// createDatabase creates the database schema if it doesn't exist
func createDatabase(config *Config) error {
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.DBName))
	return err
}

// connectToDatabase connects to the newly created database
func connectToDatabase(config *Config) error {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	var err error
	db, err = sql.Open("mysql", connStr)
	return err
}

// createTables creates the necessary tables if they do not exist
func createTables() {
	createTableQueries := []string{
		`CREATE TABLE IF NOT EXISTS users (id INT AUTO_INCREMENT PRIMARY KEY);`,
		`CREATE TABLE IF NOT EXISTS artists (id INT AUTO_INCREMENT PRIMARY KEY, name TEXT);`,
		`CREATE TABLE IF NOT EXISTS formats (id INT AUTO_INCREMENT PRIMARY KEY, name TEXT, description TEXT);`,
		`CREATE TABLE IF NOT EXISTS media (id INT AUTO_INCREMENT PRIMARY KEY, title TEXT, date_published DATE, image_url TEXT, genre_tags TEXT, artist_id INT, format_id INT, CONSTRAINT fk_media_artist FOREIGN KEY (artist_id) REFERENCES artists(id), CONSTRAINT fk_media_format FOREIGN KEY (format_id) REFERENCES formats(id), CONSTRAINT unique_media UNIQUE (title(255), artist_id, format_id));`,
		`CREATE TABLE IF NOT EXISTS user_media (user_id INT, media_id INT, format_id INT, PRIMARY KEY (user_id, media_id, format_id), CONSTRAINT fk_user_media_user FOREIGN KEY (user_id) REFERENCES users(id), CONSTRAINT fk_user_media_media FOREIGN KEY (media_id) REFERENCES media(id), CONSTRAINT fk_user_media_format FOREIGN KEY (format_id) REFERENCES formats(id));`,
	}
	for _, query := range createTableQueries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}

	checkAndAddColumn("artists", "name", "TEXT")
	checkAndAddColumn("formats", "name", "TEXT")
	checkAndAddColumn("formats", "description", "TEXT")
	checkAndAddColumn("media", "title", "TEXT")
	checkAndAddColumn("media", "date_published", "DATE")
	checkAndAddColumn("media", "image_url", "TEXT")
	checkAndAddColumn("media", "genre_tags", "TEXT")
	checkAndAddColumn("media", "artist_id", "INT")
	checkAndAddColumn("media", "format_id", "INT")
}

// checkAndAddColumn checks if a column exists and adds it if it doesn't
func checkAndAddColumn(tableName, columnName, columnType string) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_NAME='%s' AND COLUMN_NAME='%s';`, tableName, columnName)
	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		query = fmt.Sprintf(`ALTER TABLE %s ADD COLUMN %s %s;`, tableName, columnName, columnType)
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// importFormatsByFile populates the formats table with necessary format IDs
func importFormatsByFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Failed to open JSON file:", err)
	}
	defer file.Close()

	var formats []Format
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&formats)
	if err != nil {
		log.Fatal("Failed to decode JSON file:", err)
	}

	for _, format := range formats {
		var formatID int
		err := db.QueryRow(`SELECT id FROM formats WHERE name = ? AND description = ?`, format.Name, format.Description).Scan(&formatID)
		if err == sql.ErrNoRows {
			_, err = db.Exec(`INSERT INTO formats (name, description) VALUES (?, ?)`, format.Name, format.Description)
			if err != nil {
				log.Fatal("Failed to insert format:", err)
			}
		} else if err != nil {
			log.Fatal("Failed to query format:", err)
		}
	}
}

func importMediaByFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Failed to open JSON file:", err)
	}
	defer file.Close()

	var albums []Album
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&albums)
	if err != nil {
		log.Fatal("Failed to decode JSON file:", err)
	}

	for _, album := range albums {
		var artistID int
		err := db.QueryRow(`SELECT id FROM artists WHERE name = ?`, album.Artist).Scan(&artistID)
		if err == sql.ErrNoRows {
			// Artist not found, insert new artist
			result, err := db.Exec(`INSERT INTO artists (name) VALUES (?)`, album.Artist)
			if err != nil {
				log.Fatal("Failed to insert artist:", err)
			}
			artistID64, err := result.LastInsertId()
			if err != nil {
				log.Fatal("Failed to retrieve last insert ID for artist:", err)
			}
			artistID = int(artistID64)
		} else if err != nil {
			log.Fatal("Failed to query artist:", err)
		}

		var formatID int
		err = db.QueryRow(`SELECT id FROM formats WHERE name = ?`, album.Format).Scan(&formatID)
		if err == sql.ErrNoRows {
			log.Fatal("Format not found:", album.Format)
		} else if err != nil {
			log.Fatal("Failed to query format:", err)
		}

		_, err = db.Exec(`INSERT INTO media (title, date_published, image_url, genre_tags, artist_id, format_id) VALUES (?, ?, ?, ?, ?, ?)`,
			album.Title, album.DatePublished, album.ImageURL, strings.Join(album.GenreTags, ","), artistID, formatID)
		if err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				continue
			}
			log.Fatal("Failed to insert media:", err)
		}
	}
}

// initDB initializes the database connection and creates the schema
func initDB() error {
	config, err := loadConfig()

	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	err = connectToMySQL(config)

	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	err = createDatabase(config)

	if err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}

	err = connectToDatabase(config)

	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	createTables()
	importFormatsByFile("formats.json")
	importMediaByFile("media.json")

	return nil
}
