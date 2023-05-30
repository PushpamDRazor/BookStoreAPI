package lms

import (
	"Demo/config"
	"Demo/internal/entities/library"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateGenre(t *testing.T) {
	router := gin.Default()

	db := config.LoadConfig()
	db.Exec("TRUNCATE TABLE libraries;")

	router.POST("/libraries", StoreBook(db))

	bookPayload := library.Library{
		Book_ID: 1,
		Aisle: 1,
		Level: 2,
		Position: 3,
	}

	payloadJSON, _ := json.Marshal(bookPayload)

	req, _ := http.NewRequest("POST", "/libraries", bytes.NewBuffer(payloadJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status code %d but got %d", http.StatusCreated, resp.Code)
	}

	var createdBook library.Library
	err := json.Unmarshal(resp.Body.Bytes(), &createdBook)
	if err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}

	assert.Equal(t, bookPayload.Book_ID, createdBook.Book_ID)
	assert.Equal(t, bookPayload.Aisle, createdBook.Aisle)
	assert.Equal(t, bookPayload.Level, createdBook.Level)
	assert.Equal(t, bookPayload.Position, createdBook.Position)
}

func TestGetPositionByID(t *testing.T) {
	router := gin.Default()
	db := config.LoadConfig()
	db.Exec("TRUNCATE TABLE libraries;")
	router.GET("/libraries/:id", GetPositionByID(db))

	mockBook := &library.Library{
		Book_ID: 1,
		Aisle: 1,
		Level: 2,
		Position: 3,
	}

	db.Create(mockBook)

	req, _ := http.NewRequest("GET", "/libraries/1", nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, resp.Code)
	}

	var fetchedBook library.Library
	err := json.Unmarshal(resp.Body.Bytes(), &fetchedBook)
	if err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}

	assert.Equal(t, mockBook.Book_ID, fetchedBook.Book_ID)
	assert.Equal(t, mockBook.Aisle, fetchedBook.Aisle)
	assert.Equal(t, mockBook.Level, fetchedBook.Level)
	assert.Equal(t, mockBook.Position, fetchedBook.Position)
}

func TestGetBooksPositionByAuthor(t *testing.T) {
	router := gin.Default()
	db := config.LoadConfig()
	db.Exec("TRUNCATE TABLE libraries;")
	router.GET("/author/:id", GetBooksPositionByAuthor(db))

	mockBook := &library.Library{
		Book_ID: 1,
		Aisle: 1,
		Level: 2,
		Position: 3,
	}

	db.Create(mockBook)

	req, _ := http.NewRequest("GET", "/author/1", nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, resp.Code)
	}

	var fetchedBooks []library.Library
	err := json.Unmarshal(resp.Body.Bytes(), &fetchedBooks)
	if err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}

	fetchedBook := fetchedBooks[0]
	assert.Equal(t, mockBook.Book_ID, fetchedBook.Book_ID)
	assert.Equal(t, mockBook.Aisle, fetchedBook.Aisle)
	assert.Equal(t, mockBook.Level, fetchedBook.Level)
	assert.Equal(t, mockBook.Position, fetchedBook.Position)
}