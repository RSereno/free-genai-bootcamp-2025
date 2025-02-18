package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"backend_go/db" // Import your db package
	"backend_go/models"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

// dbConn is the database connection variable
var dbConn *sql.DB

// Add this function before main()
func SetupRoutes(router *gin.Engine) {
	// Move all route registrations here from main()
	router.GET("/api/ping", pingHandler)
	router.GET("/api/words", getWordsHandler)
	router.GET("/api/words/:id", getWordByIDHandler)
	router.POST("/api/words", createWordHandler)
	router.PUT("/api/words/:id", updateWordHandler)
	router.DELETE("/api/words/:id", deleteWordHandler)
	router.GET("/api/groups", getGroupsHandler)
	router.GET("/api/groups/:id", getGroupByIDHandler)
	router.POST("/api/groups", createGroupHandler)
	router.PUT("/api/groups/:id", updateGroupHandler)
	router.DELETE("/api/groups/:id", deleteGroupHandler)
	router.GET("/api/study_sessions", getStudySessionsHandler)
	router.GET("/api/study_sessions/:id", getStudySessionByIDHandler)
	router.POST("/api/study_sessions", createStudySessionHandler)
	router.PUT("/api/study_sessions/:id", updateStudySessionHandler)
	router.DELETE("/api/study_sessions/:id", deleteStudySessionHandler)
	router.GET("/api/study_activities", getStudyActivitiesHandler)
	router.GET("/api/study_activities/:id", getStudyActivityByIDHandler)
	router.POST("/api/study_activities", createStudyActivityHandler)
	router.PUT("/api/study_activities/:id", updateStudyActivityHandler)
	router.DELETE("/api/study_activities/:id", deleteStudyActivityHandler)
	router.GET("/api/word_review_items", getWordReviewItemsHandler)
	router.GET("/api/word_review_items/:id", getWordReviewItemByIDHandler)
	router.POST("/api/word_review_items", createWordReviewItemHandler)
	router.PUT("/api/word_review_items/:id", updateWordReviewItemHandler)
	router.DELETE("/api/word_review_items/:id", deleteWordReviewItemHandler)
	router.GET("/api/words_groups", getWordsGroupsHandler)
	router.GET("/api/words_groups/:id", getWordsGroupsByIDHandler)
	router.POST("/api/words_groups", createWordsGroupsHandler)
	router.PUT("/api/words_groups/:id", updateWordsGroupsHandler)
	router.DELETE("/api/words_groups/:id", deleteWordsGroupsHandler)
}

func main() {
	// Initialize database connection
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbConn.Close()

	router := gin.Default()
	SetupRoutes(router) // Now uses the shared function

	// Start the server and handle errors
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// getWordsHandler handles the /api/words endpoint.
func getWordsHandler(c *gin.Context) {
	words, err := db.GetAllWords(dbConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch words from database"})
		log.Println("Failed to fetch words:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": words})
}

// getWordByIDHandler handles the /api/words/:id endpoint.
func getWordByIDHandler(c *gin.Context) {
	idStr := c.Param("id") // Get the "id" parameter from the URL

	// Convert the id from string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	// Call the db.GetWordByID function to retrieve the word from the database
	word, err := db.GetWordByID(dbConn, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch word from database"})
		log.Println("Failed to fetch word:", err)
		return
	}

	// If the word is not found, return a 404 Not Found error
	if word == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
		return
	}

	// Return the word as a JSON response
	c.JSON(http.StatusOK, gin.H{"item": word})
}

// getGroupsHandler handles the /api/groups endpoint.
func getGroupsHandler(c *gin.Context) {
	groups, err := db.GetAllGroups(dbConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups from database"})
		log.Println("Failed to fetch groups:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": groups})
}

// getGroupByIDHandler handles the /api/groups/:id endpoint.
func getGroupByIDHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	group, err := db.GetGroupByID(dbConn, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch group from database"})
		log.Println("Failed to fetch group:", err)
		return
	}

	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": group})
}

func initDB() error {
	dbPath := filepath.Join(".", "words.db") // Path to the database file

	// Check if the database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("Database file does not exist. Please run 'mage initdb' to create it.")
		return fmt.Errorf("database file not found: %s. Run 'mage initdb'", dbPath)
	}

	fmt.Println("Connecting to database...")
	var err error
	dbConn, err = sql.Open("sqlite3", dbPath) // Open connection to SQLite database
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := dbConn.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Successfully connected to database.")
	return nil
}

// createWordHandler handles the POST /api/words endpoint.
func createWordHandler(c *gin.Context) {
	var word models.Word
	if err := c.BindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	id, err := db.CreateWord(dbConn, &word)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create word in database"})
		log.Println("Failed to create word:", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// updateWordHandler handles the PUT /api/words/:id endpoint.
func updateWordHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	var word models.Word
	if err := c.BindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	word.ID = id

	if err := db.UpdateWord(dbConn, &word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update word in database"})
		log.Println("Failed to update word:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Word updated successfully"})
}

// deleteWordHandler handles the DELETE /api/words/:id endpoint.
func deleteWordHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	// Call the db.DeleteWord function to delete the word from the database
	if err := db.DeleteWord(dbConn, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete word from database"})
		log.Println("Failed to delete word:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Word deleted successfully"})
}

// createGroupHandler handles the POST /api/groups endpoint.
func createGroupHandler(c *gin.Context) {
	var group models.Group

	// Bind the JSON data from the request body to the group struct
	if err := c.BindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the db.CreateGroup function to create the group in the database
	id, err := db.CreateGroup(dbConn, &group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group in database"})
		log.Println("Failed to create group:", err)
		return
	}

	// Return the ID of the newly created group in the response
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// updateGroupHandler handles the PUT /api/groups/:id endpoint.
func updateGroupHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var group models.Group
	if err := c.BindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	group.ID = id // Set the ID of the group to the ID from the URL

	// Call the db.UpdateGroup function to update the group in the database
	if err := db.UpdateGroup(dbConn, &group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update group in database"})
		log.Println("Failed to update group:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group updated successfully"})
}

// deleteGroupHandler handles the DELETE /api/groups/:id endpoint.
func deleteGroupHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Call the db.DeleteGroup function to delete the group from the database
	if err := db.DeleteGroup(dbConn, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group from database"})
		log.Println("Failed to delete group:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}

// getStudySessionsHandler handles the /api/study_sessions endpoint.
func getStudySessionsHandler(c *gin.Context) {
	studySessions, err := db.GetAllStudySessions(dbConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch study sessions from database"})
		log.Println("Failed to fetch study sessions:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": studySessions})
}

// getStudySessionByIDHandler handles the /api/study_sessions/:id endpoint.
func getStudySessionByIDHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid study session ID"})
		return
	}

	studySession, err := db.GetStudySessionByID(dbConn, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch study session from database"})
		log.Println("Failed to fetch study session:", err)
		return
	}

	if studySession == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Study session not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": studySession})
}

// createStudySessionHandler handles the POST /api/study_sessions endpoint.
func createStudySessionHandler(c *gin.Context) {
	var studySession models.StudySession

	// Bind the JSON data from the request body to the studySession struct
	if err := c.BindJSON(&studySession); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the db.CreateStudySession function to create the studySession in the database
	id, err := db.CreateStudySession(dbConn, &studySession)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create study session in database"})
		log.Println("Failed to create study session:", err)
		return
	}

	// Return the ID of the newly created studySession in the response
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// updateStudySessionHandler handles the PUT /api/study_sessions/:id endpoint.
func updateStudySessionHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid study session ID"})
		return
	}

	var studySession models.StudySession
	if err := c.BindJSON(&studySession); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	studySession.ID = id // Set the ID of the studySession to the ID from the URL

	// Call the db.UpdateStudySession function to update the studySession in the database
	if err := db.UpdateStudySession(dbConn, &studySession); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update study session in database"})
		log.Println("Failed to update study session:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Study session updated successfully"})
}

// deleteStudySessionHandler handles the DELETE /api/study_sessions/:id endpoint.
func deleteStudySessionHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid study session ID"})
		return
	}

	// Call the db.DeleteStudySession function to delete the studySession from the database
	if err := db.DeleteStudySession(dbConn, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete study session from database"})
		log.Println("Failed to delete study session:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Study session deleted successfully"})
}

// getStudyActivitiesHandler handles the /api/study_activities endpoint.
func getStudyActivitiesHandler(c *gin.Context) {
	studyActivities, err := db.GetAllStudyActivities(dbConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch study activities from database"})
		log.Println("Failed to fetch study activities:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": studyActivities})
}

// getStudyActivityByIDHandler handles the /api/study_activities/:id endpoint.
func getStudyActivityByIDHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid study activity ID"})
		return
	}

	studyActivity, err := db.GetStudyActivityByID(dbConn, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch study activity from database"})
		log.Println("Failed to fetch study activity:", err)
		return
	}

	if studyActivity == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Study activity not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": studyActivity})
}

// createStudyActivityHandler handles the POST /api/study_activities endpoint.
func createStudyActivityHandler(c *gin.Context) {
	var studyActivity models.StudyActivity

	// Bind the JSON data from the request body to the studyActivity struct
	if err := c.BindJSON(&studyActivity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the db.CreateStudyActivity function to create the studyActivity in the database
	id, err := db.CreateStudyActivity(dbConn, &studyActivity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create study activity in database"})
		log.Println("Failed to create study activity:", err)
		return
	}

	// Return the ID of the newly created studyActivity in the response
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// updateStudyActivityHandler handles the PUT /api/study_activities/:id endpoint.
func updateStudyActivityHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid study activity ID"})
		return
	}

	var studyActivity models.StudyActivity
	if err := c.BindJSON(&studyActivity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	studyActivity.ID = id // Set the ID of the studyActivity to the ID from the URL

	// Call the db.UpdateStudyActivity function to update the studyActivity in the database
	if err := db.UpdateStudyActivity(dbConn, &studyActivity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update study activity in database"})
		log.Println("Failed to update study activity:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Study activity updated successfully"})
}

// deleteStudyActivityHandler handles the DELETE /api/study_activities/:id endpoint.
func deleteStudyActivityHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid study activity ID"})
		return
	}

	// Call the db.DeleteStudyActivity function to delete the studyActivity from the database
	if err := db.DeleteStudyActivity(dbConn, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete study activity from database"})
		log.Println("Failed to delete study activity:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Study activity deleted successfully"})
}

// getWordReviewItemsHandler handles the /api/word_review_items endpoint.
func getWordReviewItemsHandler(c *gin.Context) {
	wordReviewItems, err := db.GetAllWordReviewItems(dbConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch word review items from database"})
		log.Println("Failed to fetch word review items:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": wordReviewItems})
}

// getWordReviewItemByIDHandler handles the /api/word_review_items/:id endpoint.
func getWordReviewItemByIDHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word review item ID"})
		return
	}

	wordReviewItem, err := db.GetWordReviewItemByID(dbConn, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch word review item from database"})
		log.Println("Failed to fetch word review item:", err)
		return
	}

	if wordReviewItem == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Word review item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": wordReviewItem})
}

// createWordReviewItemHandler handles the POST /api/word_review_items endpoint.
func createWordReviewItemHandler(c *gin.Context) {
	var wordReviewItem models.WordReviewItem

	// Bind the JSON data from the request body to the wordReviewItem struct
	if err := c.BindJSON(&wordReviewItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the db.CreateWordReviewItem function to create the wordReviewItem in the database
	id, err := db.CreateWordReviewItem(dbConn, &wordReviewItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create word review item in database"})
		log.Println("Failed to create word review item:", err)
		return
	}

	// Return the ID of the newly created wordReviewItem in the response
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// updateWordReviewItemHandler handles the PUT /api/word_review_items/:id endpoint.
func updateWordReviewItemHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word review item ID"})
		return
	}

	var wordReviewItem models.WordReviewItem
	if err := c.BindJSON(&wordReviewItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	wordReviewItem.ID = id // Set the ID of the wordReviewItem to the ID from the URL

	// Call the db.UpdateWordReviewItem function to update the wordReviewItem in the database
	if err := db.UpdateWordReviewItem(dbConn, &wordReviewItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update word review item in database"})
		log.Println("Failed to update word review item:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Word review item updated successfully"})
}

// deleteWordReviewItemHandler handles the DELETE /api/word_review_items/:id endpoint.
func deleteWordReviewItemHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word review item ID"})
		return
	}

	// Call the db.DeleteWordReviewItem function to delete the wordReviewItem from the database
	if err := db.DeleteWordReviewItem(dbConn, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete word review item from database"})
		log.Println("Failed to delete word review item:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Word review item deleted successfully"})
}

// getWordsGroupsHandler handles the /api/words_groups endpoint.
func getWordsGroupsHandler(c *gin.Context) {
	wordsGroups, err := db.GetAllWordsGroups(dbConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch words_groups from database"})
		log.Println("Failed to fetch words_groups:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": wordsGroups})
}

// getWordsGroupsByIDHandler handles the /api/words_groups/:id endpoint.
func getWordsGroupsByIDHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid words_groups ID"})
		return
	}

	wordsGroup, err := db.GetWordsGroupsByID(dbConn, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch words_groups from database"})
		log.Println("Failed to fetch words_groups:", err)
		return
	}

	if wordsGroup == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "WordsGroups not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": wordsGroup})
}

// createWordsGroupsHandler handles the POST /api/words_groups endpoint.
func createWordsGroupsHandler(c *gin.Context) {
	var wordsGroup models.WordsGroups

	// Bind the JSON data from the request body to the wordsGroup struct
	if err := c.BindJSON(&wordsGroup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the db.CreateWordsGroups function to create the wordsGroup in the database
	id, err := db.CreateWordsGroups(dbConn, &wordsGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create words_groups in database"})
		log.Println("Failed to create words_groups:", err)
		return
	}

	// Return the ID of the newly created wordsGroup in the response
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// updateWordsGroupsHandler handles the PUT /api/words_groups/:id endpoint.
func updateWordsGroupsHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid words_groups ID"})
		return
	}

	var wordsGroup models.WordsGroups
	if err := c.BindJSON(&wordsGroup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	wordsGroup.ID = id // Set the ID of the wordsGroup to the ID from the URL

	// Call the db.UpdateWordsGroups function to update the wordsGroup in the database
	if err := db.UpdateWordsGroups(dbConn, &wordsGroup); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update words_groups in database"})
		log.Println("Failed to update words_groups:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "WordsGroups updated successfully"})
}

// deleteWordsGroupsHandler handles the DELETE /api/words_groups/:id endpoint.
func deleteWordsGroupsHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid words_groups ID"})
		return
	}

	// Call the db.DeleteWordsGroups function to delete the wordsGroup from the database
	if err := db.DeleteWordsGroups(dbConn, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete words_groups from database"})
		log.Println("Failed to delete words_groups:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "WordsGroups deleted successfully"})
}
