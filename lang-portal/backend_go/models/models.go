package models

import "time"

// Word represents the 'words' table in the database.
type Word struct {
	ID         int    `json:"id"`
	English    string `json:"english"`
	Portuguese string `json:"portuguese"`
	Parts      string `json:"parts"`
}

// Group represents the 'groups' table.
type Group struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// StudySession represents the 'study_sessions' table.
type StudySession struct {
	ID              int `json:"id"`
	GroupID         int
	CreatedAt       string
	StudyActivityID int
	StartedAt       time.Time `json:"started_at"`
	EndedAt         time.Time `json:"ended_at"`
}

// StudyActivity represents the 'study_activities' table.
type StudyActivity struct {
	ID             int       `json:"id"`
	StudySessionID int       `json:"study_session_id"`
	ActivityType   string    `json:"activity_type"`
	GroupID        int       `json:"group_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// WordReviewItem represents the 'word_review_items' table.
type WordReviewItem struct {
	ID             int       `json:"id"`
	WordID         int       `json:"word_id"`
	StudySessionID int       `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

// WordsGroups represents a words_groups in the database.
type WordsGroups struct {
	ID      int `json:"id"`
	WordID  int `json:"word_id"`
	GroupID int `json:"group_id"`
}
