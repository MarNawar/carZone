package driver

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	_ "github.com/lib/pq"

)

var db *sql.DB 
func InitDB() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	log.Println(connStr)

	const maxRetries = 5          
	const retryInterval = 2 * time.Second

	var err error
	for i := 1; i <= maxRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Attempt %d: Error opening database: %v", i, err)
		} else if err = db.Ping(); err == nil {
			fmt.Println("Successfully connected to the database")
			return
		}

		log.Printf("Attempt %d: Could not connect to database, retrying in %v...", i, retryInterval)
		time.Sleep(retryInterval)
	}

	log.Fatalf("Failed to connect to the database after %d attempts: %v", maxRetries, err)
}

func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("Database connection is not initialized. Call InitDB first.")
	}
	return db
}

func CloseDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing the Database: %v", err)
		}
	}
}
