package db

import (
	"database/sql"
	"fmt"
	"log"

	"backend_go/models" // Import your models package
)

// GetAllWordReviewItems retrieves all word review items from the database.
func GetAllWordReviewItems(db *sql.DB) ([]models.WordReviewItem, error) {
	rows, err := db.Query("SELECT id, study_activity_id, word_id, is_correct FROM word_review_items")
	if err != nil {
		return nil, fmt.Errorf("failed to query word review items: %w", err)
	}
	defer rows.Close()

	var wordReviewItems []models.WordReviewItem
	for rows.Next() {
		var wordReviewItem models.WordReviewItem
		if err := rows.Scan(&wordReviewItem.ID, &wordReviewItem.StudyActivityID, &wordReviewItem.WordID, &wordReviewItem.IsCorrect); err != nil {
			log.Println("Error scanning word review item row:", err)
			continue
		}
		wordReviewItems = append(wordReviewItems, wordReviewItem)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating word review item rows: %w", err)
	}

	return wordReviewItems, nil
}

// GetWordReviewItemByID retrieves a word review item from the database by its ID.
func GetWordReviewItemByID(db *sql.DB, id int) (*models.WordReviewItem, error) {
	row := db.QueryRow("SELECT id, study_activity_id, word_id, is_correct FROM word_review_items WHERE id = ?", id)

	var wordReviewItem models.WordReviewItem
	err := row.Scan(&wordReviewItem.ID, &wordReviewItem.StudyActivityID, &wordReviewItem.WordID, &wordReviewItem.IsCorrect)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Word review item not found
		}
		return nil, fmt.Errorf("failed to scan word review item row: %w", err)
	}

	return &wordReviewItem, nil
}

// CreateWordReviewItem creates a new word review item in the database.
func CreateWordReviewItem(db *sql.DB, wordReviewItem *models.WordReviewItem) (int, error) {
	result, err := db.Exec("INSERT INTO word_review_items (study_activity_id, word_id, is_correct) VALUES (?, ?, ?)",
		wordReviewItem.StudyActivityID, wordReviewItem.WordID, wordReviewItem.IsCorrect)
	if err != nil {
		return 0, fmt.Errorf("failed to create word review item: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return int(id), nil
}

// UpdateWordReviewItem updates an existing word review item in the database.
func UpdateWordReviewItem(db *sql.DB, wordReviewItem *models.WordReviewItem) error {
	result, err := db.Exec("UPDATE word_review_items SET study_activity_id = ?, word_id = ?, is_correct = ? WHERE id = ?",
		wordReviewItem.StudyActivityID, wordReviewItem.WordID, wordReviewItem.IsCorrect, wordReviewItem.ID)
	if err != nil {
		return fmt.Errorf("failed to update word review item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("word review item with id %d not found", wordReviewItem.ID)
	}

	return nil
}

// DeleteWordReviewItem deletes a word review item from the database.
func DeleteWordReviewItem(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM word_review_items WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete word review item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("word review item with id %d not found", id)
	}

	return nil
}
