package db

import (
	"database/sql"
	"fmt"
	"log"

	"backend_go/models" // Import your models package
)

// GetAllStudyActivities retrieves all study activities from the database.
func GetAllStudyActivities(db *sql.DB) ([]models.StudyActivity, error) {
	rows, err := db.Query("SELECT id, study_session_id, group_id, created_at FROM study_activities")
	if err != nil {
		return nil, fmt.Errorf("failed to query study activities: %w", err)
	}
	defer rows.Close()

	var studyActivities []models.StudyActivity
	for rows.Next() {
		var studyActivity models.StudyActivity
		if err := rows.Scan(&studyActivity.ID, &studyActivity.StudySessionID, &studyActivity.GroupID, &studyActivity.CreatedAt); err != nil {
			log.Println("Error scanning study activity row:", err)
			continue
		}
		studyActivities = append(studyActivities, studyActivity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating study activity rows: %w", err)
	}

	return studyActivities, nil
}

// GetStudyActivityByID retrieves a study activity from the database by its ID.
func GetStudyActivityByID(db *sql.DB, id int) (*models.StudyActivity, error) {
	row := db.QueryRow("SELECT id, study_session_id, group_id, created_at FROM study_activities WHERE id = ?", id)

	var studyActivity models.StudyActivity
	err := row.Scan(&studyActivity.ID, &studyActivity.StudySessionID, &studyActivity.GroupID, &studyActivity.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Study activity not found
		}
		return nil, fmt.Errorf("failed to scan study activity row: %w", err)
	}

	return &studyActivity, nil
}

// CreateStudyActivity creates a new study activity in the database.
func CreateStudyActivity(db *sql.DB, studyActivity *models.StudyActivity) (int, error) {
	result, err := db.Exec("INSERT INTO study_activities (study_session_id, group_id, created_at) VALUES (?, ?, ?)",
		studyActivity.StudySessionID, studyActivity.GroupID, studyActivity.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to create study activity: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return int(id), nil
}

// UpdateStudyActivity updates an existing study activity in the database.
func UpdateStudyActivity(db *sql.DB, studyActivity *models.StudyActivity) error {
	result, err := db.Exec("UPDATE study_activities SET study_session_id = ?, group_id = ?, created_at = ? WHERE id = ?",
		studyActivity.StudySessionID, studyActivity.GroupID, studyActivity.CreatedAt, studyActivity.ID)
	if err != nil {
		return fmt.Errorf("failed to update study activity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("study activity with id %d not found", studyActivity.ID)
	}

	return nil
}

// DeleteStudyActivity deletes a study activity from the database.
func DeleteStudyActivity(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM study_activities WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete study activity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("study activity with id %d not found", id)
	}

	return nil
}
