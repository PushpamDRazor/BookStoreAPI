package lms

import (
	"Demo/config"
	"Demo/internal/entities/book"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateBooks(t *testing.T) {
	router := gin.Default()

	db := config.LoadConfig()
	db.Exec("TRUNCATE TABLE books;")

	router.POST("/books", CreateBooks(db))

	bookPayload := book.Book{
		Title: "Merchant of Venice",
		AuthorName: "John Doe",
		GenreName: "Horror",
		PublicationDate: "1677312000",
	}

	payloadJSON, _ := json.Marshal(bookPayload)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(payloadJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status code %d but got %d", http.StatusCreated, resp.Code)
	}

	var createdBook book.Book
	err := json.Unmarshal(resp.Body.Bytes(), &createdBook)
	if err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}

	assert.Equal(t, bookPayload.Title, createdBook.Title)
	assert.Equal(t, bookPayload.AuthorName, createdBook.AuthorName)
	assert.Equal(t, bookPayload.GenreName, createdBook.GenreName)
	assert.Equal(t, "2023-02-25", createdBook.PublicationDate)
}

func TestGetBookById(t *testing.T) {
	router := gin.Default()
	db := config.LoadConfig()
	db.Exec("TRUNCATE TABLE books;")
	router.GET("/books/:id", GetBookById(db))

	bookPayload := book.Book{
		ID: 1,
		Title: "Merchant of Venice",
		AuthorName: "John Doe",
		GenreName: "Horror",
		PublicationDate: "2023-02-25",
		CreatedAt:   "2023-05-28 12:34:56",
		UpdatedAt:   "2023-05-28 12:34:56",
		IsDeleted:   false,
	}

	db.Create(bookPayload)

	req, _ := http.NewRequest("GET", "/books/1", nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, resp.Code)
	}

	var fetchedBooks []book.Book
	err := json.Unmarshal(resp.Body.Bytes(), &fetchedBooks)
	if err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}

	fetchedBook := fetchedBooks[0]
	assert.Equal(t, bookPayload.Title, fetchedBook.Title)
	assert.Equal(t, bookPayload.AuthorName, fetchedBook.AuthorName)
	assert.Equal(t, bookPayload.GenreName, fetchedBook.GenreName)
	assert.Equal(t, "2023-02-25", fetchedBook.PublicationDate)
}