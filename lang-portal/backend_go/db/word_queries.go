package db

import (
	"database/sql"
	"fmt"
	"log"

	"backend_go/models" // Import your models package
)

// GetAllWords retrieves all words from the database.
func GetAllWords(db *sql.DB) ([]models.Word, error) {
	rows, err := db.Query("SELECT id, english, portuguese, parts FROM words")
	if err != nil {
		return nil, fmt.Errorf("failed to query words: %w", err)
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		if err := rows.Scan(&word.ID, &word.English, &word.Portuguese, &word.Parts); err != nil {
			log.Println("Error scanning word row:", err)
			continue
		}
		words = append(words, word)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating words rows: %w", err)
	}

	return words, nil
}

// GetWordByID retrieves a word from the database by its ID.
func GetWordByID(db *sql.DB, id int) (*models.Word, error) {
	row := db.QueryRow("SELECT id, english, portuguese, parts FROM words WHERE id = ?", id)

	var word models.Word
	err := row.Scan(&word.ID, &word.English, &word.Portuguese, &word.Parts)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Word not found
		}
		return nil, fmt.Errorf("failed to scan word row: %w", err)
	}

	return &word, nil
}
