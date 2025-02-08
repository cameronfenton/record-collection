package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// db is the global database connection pool
var db *sql.DB

// Config struct holds the database configuration details
type Config struct {
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
}

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
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort)
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
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	var err error
	db, err = sql.Open("mysql", connStr)
	return err
}

// initDB initializes the database connection and creates the schema
func initDB() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	err = connectToMySQL(config)
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}

	err = createDatabase(config)
	if err != nil {
		log.Fatal("Failed to create database:", err)
	}

	err = connectToDatabase(config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	createTables()
	insertDummyRecords()
}

// createTables creates the necessary tables if they do not exist
func createTables() {
	// Create tables if they do not exist
	createTableQueries := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT PRIMARY KEY
        );`,
		`CREATE TABLE IF NOT EXISTS artists (
            id INT AUTO_INCREMENT PRIMARY KEY
        );`,
		`CREATE TABLE IF NOT EXISTS formats (
            id INT AUTO_INCREMENT PRIMARY KEY
        );`,
		`CREATE TABLE IF NOT EXISTS media (
            id INT AUTO_INCREMENT PRIMARY KEY
        );`,
		`CREATE TABLE IF NOT EXISTS user_media (
            user_id INT,
            media_id INT,
            format_id INT,
            PRIMARY KEY (user_id, media_id, format_id)
        );`,
	}

	// Run create table queries
	for _, query := range createTableQueries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Add columns if they do not exist
	checkAndAddColumn("users", "username", "VARCHAR(100) NOT NULL")
	checkAndAddColumn("users", "email", "VARCHAR(100) NOT NULL")
	checkAndAddColumn("artists", "name", "TEXT")
	checkAndAddColumn("artists", "band_ids", "TEXT")
	checkAndAddColumn("artists", "first_name", "TEXT")
	checkAndAddColumn("artists", "last_name", "TEXT")
	checkAndAddColumn("artists", "image_url", "TEXT")
	checkAndAddColumn("artists", "is_group", "BOOLEAN DEFAULT FALSE")
	checkAndAddColumn("formats", "name", "TEXT")
	checkAndAddColumn("formats", "description", "TEXT")
	checkAndAddColumn("media", "title", "TEXT")
	checkAndAddColumn("media", "date_published", "DATE")
	checkAndAddColumn("media", "image_url", "TEXT")
	checkAndAddColumn("media", "genre_tags", "TEXT")
	checkAndAddColumn("media", "artist_id", "INT")
	checkAndAddColumn("media", "format_id", "INT")
	addForeignKey("user_media", "fk_user_media_user", "user_id", "users", "id")
	addForeignKey("user_media", "fk_user_media_media", "media_id", "media", "id")
	addForeignKey("user_media", "fk_user_media_format", "format_id", "formats", "id")
}

// checkAndAddColumn checks if a column exists and adds it if it doesn't
func checkAndAddColumn(tableName, columnName, columnType string) {
	query := fmt.Sprintf(`
        SELECT COUNT(*) 
        FROM information_schema.COLUMNS 
        WHERE TABLE_NAME='%s' AND COLUMN_NAME='%s';`, tableName, columnName)
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

// addForeignKey checks if a foreign key constraint exists and adds it if it doesn't
func addForeignKey(tableName, constraintName, columnName, refTableName, refColumnName string) {
	query := fmt.Sprintf(`
        SELECT COUNT(*) 
        FROM information_schema.TABLE_CONSTRAINTS 
        WHERE CONSTRAINT_NAME='%s' AND TABLE_NAME='%s';`, constraintName, tableName)
	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		query = fmt.Sprintf(`ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s);`,
			tableName, constraintName, columnName, refTableName, refColumnName)
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// insertDummyRecords inserts some dummy records into the database for testing
func insertDummyRecords() {
	// Insert dummy user
	_, err := db.Exec(`INSERT IGNORE INTO users (username, email) VALUES ('Cameron Fenton', 'cameron.fenton@example.com');`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert dummy artists
	artists := []struct {
		Name      string
		BandIds   string
		FirstName string
		LastName  string
		ImageUrl  string
		IsGroup   bool
	}{
		{Name: "Artist 1", BandIds: "", FirstName: "", LastName: "", ImageUrl: "", IsGroup: false},
		{Name: "Artist 2", BandIds: "", FirstName: "", LastName: "", ImageUrl: "", IsGroup: true},
		{Name: "Artist 3", BandIds: "", FirstName: "", LastName: "", ImageUrl: "", IsGroup: false},
	}

	// Check if each artist already exists before inserting
	for _, artist := range artists {
		var count int
		query := `SELECT COUNT(*) FROM artists WHERE name=? AND band_ids=? AND first_name=? AND last_name=? AND image_url=? AND is_group=?;`
		err := db.QueryRow(query, artist.Name, artist.BandIds, artist.FirstName, artist.LastName, artist.ImageUrl, artist.IsGroup).Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
		if count == 0 {
			_, err := db.Exec(`INSERT INTO artists (name, band_ids, first_name, last_name, image_url, is_group) VALUES (?, ?, ?, ?, ?, ?);`,
				artist.Name, artist.BandIds, artist.FirstName, artist.LastName, artist.ImageUrl, artist.IsGroup)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Insert dummy formats
	formats := []string{
		"('Format 1', 'Description 1')",
		"('Format 2', 'Description 2')",
		"('Format 3', 'Description 3')",
	}
	for _, format := range formats {
		_, err := db.Exec(fmt.Sprintf(`INSERT IGNORE INTO formats (name, description) VALUES %s;`, format))
		if err != nil {
			log.Fatal(err)
		}
	}

	// Insert dummy media
	media := []string{
		"('Media 1', '2023-01-01', '', '', 1, 1)",
		"('Media 2', '2023-01-02', '', '', 2, 2)",
		"('Media 3', '2023-01-03', '', '', 3, 3)",
	}
	for _, m := range media {
		_, err := db.Exec(fmt.Sprintf(`INSERT IGNORE INTO media (title, date_published, image_url, genre_tags, artist_id, format_id) VALUES %s;`, m))
		if err != nil {
			log.Fatal(err)
		}
	}

	// Insert dummy user_media
	userMedia := []string{
		"(1, 1, 1)",
		"(1, 2, 2)",
		"(1, 3, 3)",
	}
	for _, um := range userMedia {
		_, err := db.Exec(fmt.Sprintf(`INSERT IGNORE INTO user_media (user_id, media_id, format_id) VALUES %s;`, um))
		if err != nil {
			log.Fatal(err)
		}
	}
}
