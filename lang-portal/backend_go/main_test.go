package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend_go/models"
	"backend_go/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWordCRUD(t *testing.T) {
	// Setup
	db, err := testutils.SetupTestDB()
	require.NoError(t, err)
	defer db.Close()

	router := gin.Default()
	SetupRoutes(router) // You'll need to extract route setup to a separate function

	t.Run("Create and retrieve word", func(t *testing.T) {
		// Create word
		newWord := models.Word{
			English:    "hello",
			Portuguese: "olá",
			Parts:      "interjection",
		}

		body, _ := json.Marshal(newWord)
		req, _ := http.NewRequest("POST", "/api/words", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		var createResponse struct{ ID int }
		json.Unmarshal(resp.Body.Bytes(), &createResponse)
		require.NotZero(t, createResponse.ID)

		// Get word
		req, _ = http.NewRequest("GET", "/api/words/"+string(createResponse.ID), nil)
		resp = httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var word models.Word
		json.Unmarshal(resp.Body.Bytes(), &word)
		assert.Equal(t, "hello", word.English)
		assert.Equal(t, "olá", word.Portuguese)
	})

	t.Run("Update word", func(t *testing.T) {
		// Similar structure for update test
	})

	t.Run("Delete word", func(t *testing.T) {
		// Similar structure for delete test
	})
}
