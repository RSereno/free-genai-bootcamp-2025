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
		return nil, fmt.Errorf("error iterating word rows: %w", err)
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

// CreateWord creates a new word in the database.
func CreateWord(db *sql.DB, word *models.Word) (int, error) {
	result, err := db.Exec("INSERT INTO words (english, portuguese, parts) VALUES (?, ?, ?)",
		word.English, word.Portuguese, word.Parts)
	if err != nil {
		return 0, fmt.Errorf("failed to create word: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return int(id), nil
}

// UpdateWord updates an existing word in the database.
func UpdateWord(db *sql.DB, word *models.Word) error {
	result, err := db.Exec("UPDATE words SET english = ?, portuguese = ?, parts = ? WHERE id = ?",
		word.English, word.Portuguese, word.Parts, word.ID)
	if err != nil {
		return fmt.Errorf("failed to update word: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("word with id %d not found", word.ID)
	}

	return nil
}

// DeleteWord deletes a word from the database.
func DeleteWord(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM words WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete word: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("word with id %d not found", id)
	}

	return nil
}
