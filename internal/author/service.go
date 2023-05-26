package author

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateAuthor(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var author Author
		err := c.ShouldBindJSON(&author)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		result := db.Create(&author)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a new author"})
			return
		}

		c.JSON(http.StatusCreated, author)
	}
}

func GetAuthor(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var authors []Author
		result := db.Find(&authors)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
			return
		}
		c.JSON(http.StatusOK, authors)
	}
}

func UpdateAuthor(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
			return
		}

		var author Author
		result := db.First(&author, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}

		err = c.ShouldBindJSON(&author)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		result = db.Save(&author)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the book"})
			return
		}

		c.JSON(http.StatusOK, author)
	}
}

func GetAuthorParams(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		params := c.Request.URL.Query()
		paramType := params["paramType"][0]
		paramValue := params["paramValue"][0]

		var authors []Author
		var result *gorm.DB

		switch paramType {
		case "genre":
			genreID, err := strconv.Atoi(paramValue)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
				return
			}

			result = db.Table("authors").
				Joins("INNER JOIN books ON authors.id = books.author_id").
				Joins("INNER JOIN genres ON genres.id = books.genre_id").
				Where("genres.id = ?", genreID).
				Find(&authors)

		case "nationality":
			result = db.Table("authors").Where("nationality = ?", paramValue).Find(&authors)

		case "name":
			result = db.Where("author_name LIKE ?", "%"+paramValue+"%").Find(&authors)

		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid paramType"})
			return
		}

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query authors"})
			return
		}

		c.JSON(http.StatusOK, authors)
	}
}

func DeleteAuthor(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
			return
		}

		var author Author
		result := db.Find(&author, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
			return
		}

		result = db.Delete(&author)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the author"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Author deleted"})
	}
}