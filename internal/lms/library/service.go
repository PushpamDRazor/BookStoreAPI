package lms

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"Demo/internal/entities/book"
	"Demo/internal/entities/library"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StoreBook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var lib library.Library
		err := c.ShouldBindJSON(&lib)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		if lib.Book_ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book ID is required"})
			return
		}
		var books []book.Book
		var existingBook library.Library
		result := db.Where("aisle = ? AND level = ? AND position = ?", lib.Aisle, lib.Level, lib.Position).First(&existingBook)
		if result.Error == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Position already occupied"})
			return
		}
		result = db.Table("books").Where("books.id = ?", lib.Book_ID).Where("books.is_deleted = ?", false).First(&books)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book is not present in the store"})
			return
		}

		now := time.Now()
		lib.CreatedAt = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second())

		lib.UpdatedAt = lib.CreatedAt

		result = db.Create(&lib)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store the book"})
			return
		}

		c.JSON(http.StatusCreated, lib)
	}
}

func GetPositionByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookID := c.Param("id")

		var library library.Library
		if err := db.Where("book_id = ?", bookID).First(&library).Error; err != nil || library.IsDeleted {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		book_id := library.Book_ID
		position := library.Position
		aisle := library.Aisle
		level := library.Level

		c.JSON(http.StatusOK, gin.H{"book_id": book_id, "aisle": aisle, "level": level, "position": position})
	}
}

func RemoveBook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var library library.Library
		err := c.ShouldBindJSON(&library)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		result := db.Table("library").Where("aisle = ? AND level = ? AND position = ?", library.Aisle, library.Level, library.Position).Update("is_deleted", true)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found at the specified position"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Book removed successfully"})
	}
}

func GetBooksPositionByAuthor(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
			return
		}

		var book []book.Book
		result := db.Table("books").Where("author_id = ?", authorID).Find(&book)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query books"})
			return
		}

		var bookPositions []library.Library
		for _, book := range book {
			var lib library.Library
			result = db.Table("libraries").Where("book_id = ?", book.ID).First(&lib)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book position"})
				return
			}

			bookPos := library.Library{
				Book_ID:       book.ID,
				Aisle:    lib.Aisle,
				Level:    lib.Level,
				Position: lib.Position,
			}

			bookPositions = append(bookPositions, bookPos)
		}

		c.JSON(http.StatusOK, bookPositions)
	}
}

func UpdatePositionById(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Genre ID"})
			return
		}

		var lib library.Library
		result := db.First(&lib, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found in library"})
			return
		}

		err = c.ShouldBindJSON(&lib)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		now := time.Now()
		lib.UpdatedAt = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second())

		result = db.Save(&lib)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the library"})
			return
		}

		c.JSON(http.StatusOK, lib)
	}
}
