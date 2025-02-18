package db

import (
	"database/sql"
	"fmt"
	"log"

	"backend_go/models" // Import your models package
)

// GetAllWordsGroups retrieves all words_groups from the database.
func GetAllWordsGroups(db *sql.DB) ([]models.WordsGroups, error) {
	rows, err := db.Query("SELECT id, word_id, group_id FROM words_groups")
	if err != nil {
		return nil, fmt.Errorf("failed to query words_groups: %w", err)
	}
	defer rows.Close()

	var wordsGroups []models.WordsGroups
	for rows.Next() {
		var wordsGroup models.WordsGroups
		if err := rows.Scan(&wordsGroup.ID, &wordsGroup.WordID, &wordsGroup.GroupID); err != nil {
			log.Println("Error scanning words_groups row:", err)
			continue
		}
		wordsGroups = append(wordsGroups, wordsGroup)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating words_groups rows: %w", err)
	}

	return wordsGroups, nil
}

// GetWordsGroupsByID retrieves a words_groups from the database by its ID.
func GetWordsGroupsByID(db *sql.DB, id int) (*models.WordsGroups, error) {
	row := db.QueryRow("SELECT id, word_id, group_id FROM words_groups WHERE id = ?", id)

	var wordsGroup models.WordsGroups
	err := row.Scan(&wordsGroup.ID, &wordsGroup.WordID, &wordsGroup.GroupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // WordsGroups not found
		}
		return nil, fmt.Errorf("failed to scan words_groups row: %w", err)
	}

	return &wordsGroup, nil
}

// CreateWordsGroups creates a new words_groups in the database.
func CreateWordsGroups(db *sql.DB, wordsGroup *models.WordsGroups) (int, error) {
	result, err := db.Exec("INSERT INTO words_groups (word_id, group_id) VALUES (?, ?)",
		wordsGroup.WordID, wordsGroup.GroupID)
	if err != nil {
		return 0, fmt.Errorf("failed to create words_groups: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return int(id), nil
}

// UpdateWordsGroups updates an existing words_groups in the database.
func UpdateWordsGroups(db *sql.DB, wordsGroup *models.WordsGroups) error {
	result, err := db.Exec("UPDATE words_groups SET word_id = ?, group_id = ? WHERE id = ?",
		wordsGroup.WordID, wordsGroup.GroupID, wordsGroup.ID)
	if err != nil {
		return fmt.Errorf("failed to update words_groups: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("words_groups with id %d not found", wordsGroup.ID)
	}

	return nil
}

// DeleteWordsGroups deletes a words_groups from the database.
func DeleteWordsGroups(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM words_groups WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete words_groups: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("words_groups with id %d not found", id)
	}

	return nil
}
