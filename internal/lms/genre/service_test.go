package lms

import (
	"Demo/config"
	"Demo/internal/entities/genre"
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
	db.Exec("TRUNCATE TABLE genres;")

	router.POST("/genres", CreateGenre(db))

	genrePayload := genre.Genre{
		Genre : "Horror",
	}

	payloadJSON, _ := json.Marshal(genrePayload)

	req, _ := http.NewRequest("POST", "/genres", bytes.NewBuffer(payloadJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status code %d but got %d", http.StatusCreated, resp.Code)
	}

	var createdGenre genre.Genre
	err := json.Unmarshal(resp.Body.Bytes(), &createdGenre)
	if err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}

	assert.Equal(t, genrePayload.Genre, createdGenre.Genre)
}

func TestGenreById(t *testing.T) {
	router := gin.Default()
	db := config.LoadConfig()
	db.Exec("TRUNCATE TABLE genres;")
	router.GET("/genres/:id", GetGenreById(db))

	mockGenre := &genre.Genre{
		ID:          1,
		Genre : "Horror",
	}

	db.Create(mockGenre)

	req, _ := http.NewRequest("GET", "/genres/1", nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, resp.Code)
	}

	var fetchedGenres []genre.Genre
	err := json.Unmarshal(resp.Body.Bytes(), &fetchedGenres)
	if err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}

	fetchedGenre := fetchedGenres[0]
	assert.Equal(t, 1, fetchedGenre.ID)
	assert.Equal(t, mockGenre.Genre, fetchedGenre.Genre)
}