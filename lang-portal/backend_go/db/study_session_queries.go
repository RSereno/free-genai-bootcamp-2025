package db

import (
	"database/sql"
	"fmt"
	"log"

	"backend_go/models" // Import your models package
)

// GetAllStudySessions retrieves all study sessions from the database.
func GetAllStudySessions(db *sql.DB) ([]models.StudySession, error) {
	rows, err := db.Query("SELECT id, group_id, created_at, study_activity_id FROM study_sessions")
	if err != nil {
		return nil, fmt.Errorf("failed to query study sessions: %w", err)
	}
	defer rows.Close()

	var studySessions []models.StudySession
	for rows.Next() {
		var studySession models.StudySession
		if err := rows.Scan(&studySession.ID, &studySession.GroupID, &studySession.CreatedAt, &studySession.StudyActivityID); err != nil {
			log.Println("Error scanning study session row:", err)
			continue
		}
		studySessions = append(studySessions, studySession)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating study session rows: %w", err)
	}

	return studySessions, nil
}

// GetStudySessionByID retrieves a study session from the database by its ID.
func GetStudySessionByID(db *sql.DB, id int) (*models.StudySession, error) {
	row := db.QueryRow("SELECT id, group_id, created_at, study_activity_id FROM study_sessions WHERE id = ?", id)

	var studySession models.StudySession
	err := row.Scan(&studySession.ID, &studySession.GroupID, &studySession.CreatedAt, &studySession.StudyActivityID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Study session not found
		}
		return nil, fmt.Errorf("failed to scan study session row: %w", err)
	}

	return &studySession, nil
}

// CreateStudySession creates a new study session in the database.
func CreateStudySession(db *sql.DB, studySession *models.StudySession) (int, error) {
	result, err := db.Exec("INSERT INTO study_sessions (group_id, created_at, study_activity_id) VALUES (?, ?, ?)",
		studySession.GroupID, studySession.CreatedAt, studySession.StudyActivityID)
	if err != nil {
		return 0, fmt.Errorf("failed to create study session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return int(id), nil
}

// UpdateStudySession updates an existing study session in the database.
func UpdateStudySession(db *sql.DB, studySession *models.StudySession) error {
	result, err := db.Exec("UPDATE study_sessions SET group_id = ?, created_at = ?, study_activity_id = ? WHERE id = ?",
		studySession.GroupID, studySession.CreatedAt, studySession.StudyActivityID, studySession.ID)
	if err != nil {
		return fmt.Errorf("failed to update study session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("study session with id %d not found", studySession.ID)
	}

	return nil
}

// DeleteStudySession deletes a study session from the database.
func DeleteStudySession(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM study_sessions WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete study session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("study session with id %d not found", id)
	}

	return nil
}
