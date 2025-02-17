package models

import "time"

// Word represents the 'words' table in the database.
type Word struct {
	ID         int    `json:"id"`
	English    string `json:"english"`
	Portuguese string `json:"portuguese"`
	Parts      string `json:"parts"` // e.g., "verb", "noun", "adjective"
}

// WordsGroup represents the 'words_groups' join table.
type WordsGroup struct {
	ID      int `json:"id"`
	WordID  int `json:"word_id"`
	GroupID int `json:"group_id"`
}

// Group represents the 'groups' table.
type Group struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// StudySession represents the 'study_sessions' table.
type StudySession struct {
	ID              int       `json:"id"`
	GroupID         int       `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
	StudyActivityID int       `json:"study_activity_id"`
}

// StudyActivity represents the 'study_activities' table.
type StudyActivity struct {
	ID             int       `json:"id"`
	StudySessionID int       `json:"study_session_id"`
	GroupID        int       `json:"group_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// WordReviewItem represents the 'word_review_items' table.
type WordReviewItem struct {
	ID             int       `json:"id"`
	WordID         int       `json:"word_id"`
	StudySessionID int       `json:"study_session_id"`
	IsCorrect      bool      `json:"is_correct"`
	CreatedAt      time.Time `json:"created_at"`
}

/*

**Explanation of `models/word.go`:**

*   **`package models`**:  Declares that these structs belong to the `models` package.
*   **`import "time"`**: Imports the `time` package because some of our models have `CreatedAt` fields of type `time.Time`.
*   **Struct Definitions (e.g., `Word`, `Group`, etc.)**:
    *   For each table in your database schema, we've created a corresponding Go struct.
    *   **Fields:** Each struct has fields that match the columns in the database table.
    *   **Data Types:** We've chosen appropriate Go data types for each field:
        *   `int`: For integer IDs and foreign keys.
        *   `string`: For text fields like `English`, `Portuguese`, `Name`, `Description`, `Parts`.
        *   `bool`: For boolean fields like `IsCorrect`.
        *   `time.Time`: For datetime fields like `CreatedAt`.
    *   **JSON Tags (`json:"..."`)**:  Each struct field has a `json:"..."` tag. This tag is important for:
        *   **JSON Serialization:** When we want to return data from our API as JSON, the `json` tags tell the `encoding/json` package how to map Go struct fields to JSON keys. For example, the `ID` field in the `Word` struct will be serialized as `"id"` in the JSON output.

**How to use these models:**

These structs will be used throughout your backend code to:

*   **Represent data retrieved from the database:** When you query the database, you can scan the results into instances of these structs.
*   **Prepare data to be inserted into the database:** You can create instances of these structs and then insert the data into the corresponding tables.
*   **Serialize data to JSON responses:** When your API endpoints return data, you will likely be working with these structs and using Gin to automatically serialize them into JSON responses.

**Next Steps:**

1.  **Database Queries (Data Access Layer):** We can now start creating functions to interact with the database. For example, we can create functions to:
    *   Fetch a list of words.
    *   Get a specific word by ID.
    *   Fetch groups, study sessions, etc.
    *   Insert new words, groups, etc.

    We can create a new file (e.g., in a `db` subdirectory or directly in `models`, depending on project structure preference) to hold these database query functions.

Are you ready to move on to creating database query functions (Data Access Layer)? Let me know if you have any questions about the data models we just defined!

*/
