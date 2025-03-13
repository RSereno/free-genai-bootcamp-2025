package db

import (
	"database/sql"
	"fmt"

	"backend_go/models"
)

// GetStudySessionWords retrieves words reviewed in a study session with pagination
func GetStudySessionWords(db *sql.DB, sessionID, page, limit int) ([]models.Word, int, error) {
	offset := (page - 1) * limit

	// First, get total count
	var totalItems int
	countQuery := `
        SELECT COUNT(DISTINCT w.id)
        FROM words w
        JOIN word_review_items wri ON w.id = wri.word_id
        WHERE wri.study_session_id = ?
    `
	err := db.QueryRow(countQuery, sessionID).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting study session words: %v", err)
	}

	// Then get paginated words
	query := `
        SELECT DISTINCT w.*
        FROM words w
        JOIN word_review_items wri ON w.id = wri.word_id
        WHERE wri.study_session_id = ?
        ORDER BY w.id
        LIMIT ? OFFSET ?
    `
	rows, err := db.Query(query, sessionID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying study session words: %v", err)
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		err := rows.Scan(&word.ID, &word.English, &word.Portuguese, &word.Parts)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning word row: %v", err)
		}
		words = append(words, word)
	}

	return words, totalItems, nil
}

// GetStudySessionWordsRaw retrieves raw word review data for a study session with pagination
func GetStudySessionWordsRaw(db *sql.DB, sessionID, page, limit int) ([]models.WordReviewItem, int, error) {
	offset := (page - 1) * limit

	// First, get total count
	var totalItems int
	countQuery := `
        SELECT COUNT(*)
        FROM word_review_items
        WHERE study_session_id = ?
    `
	err := db.QueryRow(countQuery, sessionID).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting word review items: %v", err)
	}

	// Then get paginated word review items
	query := `
        SELECT id, word_id, study_session_id, correct, created_at
        FROM word_review_items
        WHERE study_session_id = ?
        ORDER BY created_at
        LIMIT ? OFFSET ?
    `
	rows, err := db.Query(query, sessionID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying word review items: %v", err)
	}
	defer rows.Close()

	var reviews []models.WordReviewItem
	for rows.Next() {
		var review models.WordReviewItem
		err := rows.Scan(&review.ID, &review.WordID, &review.StudySessionID, &review.Correct, &review.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning word review item row: %v", err)
		}
		reviews = append(reviews, review)
	}

	return reviews, totalItems, nil
}

// GetWordGroupStudySessions retrieves study sessions for a word group with pagination
func GetWordGroupStudySessions(db *sql.DB, groupID, page, limit int) ([]models.StudySession, int, error) {
	offset := (page - 1) * limit

	// First, get total count
	var totalItems int
	countQuery := `
        SELECT COUNT(DISTINCT ss.id)
        FROM study_sessions ss
        JOIN word_review_items wri ON ss.id = wri.study_session_id
        JOIN words_groups wg ON wg.word_id = wri.word_id
        WHERE wg.group_id = ?
    `
	err := db.QueryRow(countQuery, groupID).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting group study sessions: %v", err)
	}

	// Then get paginated study sessions
	query := `
        SELECT DISTINCT ss.*
        FROM study_sessions ss
        JOIN word_review_items wri ON ss.id = wri.study_session_id
        JOIN words_groups wg ON wg.word_id = wri.word_id
        WHERE wg.group_id = ?
        ORDER BY ss.started_at DESC
        LIMIT ? OFFSET ?
    `
	rows, err := db.Query(query, groupID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying group study sessions: %v", err)
	}
	defer rows.Close()

	var sessions []models.StudySession
	for rows.Next() {
		var session models.StudySession
		err := rows.Scan(&session.ID, &session.StartedAt, &session.EndedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning study session row: %v", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, totalItems, nil
}

// GetWordGroupStudySessionsRaw retrieves raw study session data for a word group with pagination
func GetWordGroupStudySessionsRaw(db *sql.DB, groupID, page, limit int) ([]models.StudyActivity, int, error) {
	offset := (page - 1) * limit

	// First, get total count
	var totalItems int
	countQuery := `
        SELECT COUNT(*)
        FROM study_activities sa
        JOIN word_review_items wri ON sa.study_session_id = wri.study_session_id
        JOIN words_groups wg ON wg.word_id = wri.word_id
        WHERE wg.group_id = ?
    `
	err := db.QueryRow(countQuery, groupID).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting group study activities: %v", err)
	}

	// Then get paginated study activities
	query := `
        SELECT sa.*
        FROM study_activities sa
        JOIN word_review_items wri ON sa.study_session_id = wri.study_session_id
        JOIN words_groups wg ON wg.word_id = wri.word_id
        WHERE wg.group_id = ?
        ORDER BY sa.created_at
        LIMIT ? OFFSET ?
    `
	rows, err := db.Query(query, groupID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying group study activities: %v", err)
	}
	defer rows.Close()

	var activities []models.StudyActivity
	for rows.Next() {
		var activity models.StudyActivity
		err := rows.Scan(&activity.ID, &activity.StudySessionID, &activity.ActivityType, &activity.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning study activity row: %v", err)
		}
		activities = append(activities, activity)
	}

	return activities, totalItems, nil
}
