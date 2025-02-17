package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"backend_go/db" // Import your db package

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

// dbConn is the database connection variable
var dbConn *sql.DB

func main() {
	// Initialize database connection
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbConn.Close()

	router := gin.Default()

	router.GET("/api/ping", pingHandler)
	router.GET("/api/words", getWordsHandler)        // New endpoint for getting words
	router.GET("/api/words/:id", getWordByIDHandler) // New endpoint for getting a word by ID

	router.Run(":8080") // Listen and serve on 0.0.0.0:8080
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// getWordsHandler handles the /api/words endpoint.
func getWordsHandler(c *gin.Context) {
	words, err := db.GetAllWords(dbConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch words from database"})
		log.Println("Failed to fetch words:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": words})
}

// getWordByIDHandler handles the /api/words/:id endpoint.
func getWordByIDHandler(c *gin.Context) {
	idStr := c.Param("id") // Get the "id" parameter from the URL

	// Convert the id from string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	// Call the db.GetWordByID function to retrieve the word from the database
	word, err := db.GetWordByID(dbConn, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch word from database"})
		log.Println("Failed to fetch word:", err)
		return
	}

	// If the word is not found, return a 404 Not Found error
	if word == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
		return
	}

	// Return the word as a JSON response
	c.JSON(http.StatusOK, gin.H{"item": word})
}

func initDB() error {
	dbPath := filepath.Join(".", "words.db") // Path to the database file

	// Check if the database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("Database file does not exist. Please run 'mage initdb' to create it.")
		return fmt.Errorf("database file not found: %s. Run 'mage initdb'", dbPath)
	}

	fmt.Println("Connecting to database...")
	var err error
	dbConn, err = sql.Open("sqlite3", dbPath) // Open connection to SQLite database
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := dbConn.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Successfully connected to database.")
	return nil
}
