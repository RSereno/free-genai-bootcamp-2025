package db

import (
	"database/sql"
	"fmt"
	"log"

	"backend_go/models" // Import your models package
)

// GetAllGroups retrieves all groups from the database.
func GetAllGroups(db *sql.DB) ([]models.Group, error) {
	rows, err := db.Query("SELECT id, name, description FROM groups")
	if err != nil {
		return nil, fmt.Errorf("failed to query groups: %w", err)
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.Description); err != nil {
			log.Println("Error scanning group row:", err)
			continue
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating group rows: %w", err)
	}

	return groups, nil
}

// GetGroupByID retrieves a group from the database by its ID.
func GetGroupByID(db *sql.DB, id int) (*models.Group, error) {
	row := db.QueryRow("SELECT id, name, description FROM groups WHERE id = ?", id)

	var group models.Group
	err := row.Scan(&group.ID, &group.Name, &group.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Group not found
		}
		return nil, fmt.Errorf("failed to scan group row: %w", err)
	}

	return &group, nil
}

// CreateGroup creates a new group in the database.
func CreateGroup(db *sql.DB, group *models.Group) (int, error) {
	result, err := db.Exec("INSERT INTO groups (name, description) VALUES (?, ?)",
		group.Name, group.Description)
	if err != nil {
		return 0, fmt.Errorf("failed to create group: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return int(id), nil
}

// UpdateGroup updates an existing group in the database.
func UpdateGroup(db *sql.DB, group *models.Group) error {
	result, err := db.Exec("UPDATE groups SET name = ?, description = ? WHERE id = ?",
		group.Name, group.Description, group.ID)
	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group with id %d not found", group.ID)
	}

	return nil
}

// DeleteGroup deletes a group from the database.
func DeleteGroup(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM groups WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group with id %d not found", id)
	}

	return nil
}
