package lms

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"Demo/internal/entities/author"
	"Demo/internal/entities/book"
	"Demo/internal/entities/genre"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetBookById(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Book ID"})
			return
		}

		var book []book.Book
		result := db.First(&book, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusOK, book)
	}
}

func GetBookParams(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		offsetStr := c.DefaultQuery("offset", "0")
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset value"})
			return
		}

		limit := 3
		fmt.Printf("Offset: %d\n", offset)
		fmt.Printf("Limit: %d\n", limit)

		params := c.Request.URL.Query()

		var books []book.Book
		var result *gorm.DB
		if len(params) == 0 {
			result = db.Where("books.is_deleted = ?", false).Limit(limit).Offset(offset).Find(&books)
			if result.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Books not found"})
				return
			}
		} else {
			for paramType, j := range params {
				switch paramType {
				case "genre_name":
					for _, paramValue := range j {
						genreID, err := strconv.Atoi(paramValue)
						if err != nil {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
							return
						}

						result = db.Table("books").
							Where("books.genre_id = ?", genreID).Limit(limit).Offset(offset).
							Find(&books)

					}
				case "author_name":
					for _, paramValue := range j {
						result = db.Where("author_name = ?", paramValue).Limit(limit).Offset(offset).
							Find(&books)
					}
				case "offset":
					continue
				default:
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid paramType"})
					return

				}

			}
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query books"})
				return
			}
		}

		nextOffset := offset + limit

		c.JSON(http.StatusOK, gin.H{
			"books":       books,
			"next_offset": nextOffset,
		})
	}
}

func CreateBooks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var boook book.Book
		err := c.ShouldBindJSON(&boook)
		if err != nil || boook.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		genrename := boook.GenreName
		authorname := boook.AuthorName

		now := time.Now()
		boook.CreatedAt = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second())

		boook.UpdatedAt = boook.CreatedAt

		var author author.Author
		result := db.Table("authors").
			Where("author_name = ?", authorname).
			First(&author)
		if result.Error != nil || author.IsDeleted {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve author"})
			return
		}

		boook.AuthorId = author.ID

		parsedInteger, err := strconv.ParseInt(boook.PublicationDate, 10, 64)
		if err != nil {
			fmt.Println("Error parsing integer:", err)
			return
		}

		epochTimeSeconds := parsedInteger

		epochTime := time.Unix(epochTimeSeconds, 0)

		epochDateString := epochTime.Format("2006-01-02")
		boook.PublicationDate = epochDateString

		result = db.Table("books").Where("title = ? AND author_name = ? AND genre_name = ? AND publication_date = ?", boook.Title, boook.AuthorName, boook.GenreName, boook.PublicationDate).First(&boook)
		if result.Error == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book already created"})
			return
		}

		var genre genre.Genre
		result = db.Table("genres").
			Where("genre = ?", genrename).
			First(&genre)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve genre"})
			return
		}

		boook.GenreId = genre.ID

		result = db.Create(&boook)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a new book"})
			return
		}

		c.JSON(http.StatusCreated, boook)
	}
}

func UpdateBooks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
			return
		}

		var book book.Book
		result := db.First(&book, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}

		err = c.ShouldBindJSON(&book)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		now := time.Now()
		book.UpdatedAt = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second())

		result = db.Save(&book)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the book"})
			return
		}

		c.JSON(http.StatusOK, book)
	}
}

func DeleteBook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
			return
		}

		var book book.Book
		result := db.First(&book, id)
		if result.Error != nil || book.IsDeleted {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}

		book.IsDeleted = true
		result = db.Save(&book)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the book"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
	}
}
