package db

import (
	"testing"

	"backend_go/models"
	"backend_go/testutils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateWord(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("INSERT INTO words").
		WithArgs("test", "teste", "noun").
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := CreateWord(db, &models.Word{
		English:    "test",
		Portuguese: "teste",
		Parts:      "noun",
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllWordsIntegration(t *testing.T) {
	db, err := testutils.SetupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// Insert test data
	_, err = db.Exec(`INSERT INTO words (english, portuguese, parts) VALUES 
		('hello', 'olá', 'interjection'),
		('goodbye', 'adeus', 'interjection')`)
	require.NoError(t, err)

	words, err := GetAllWords(db)
	require.NoError(t, err)
	assert.Len(t, words, 2)
}
