package router

import (
	lmsAuthor "Demo/internal/lms/author"
	lmsBook "Demo/internal/lms/book"
	lmsGenre "Demo/internal/lms/genre"
	lmsLibrary "Demo/internal/lms/library"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	r.GET("/books/:id", lmsBook.GetBookById(db))
	r.GET("/books", lmsBook.GetBookParams(db))
	r.POST("/books", lmsBook.CreateBooks(db))
	r.PATCH("/books/:id", lmsBook.UpdateBooks(db))
	r.DELETE("/books/:id", lmsBook.DeleteBook(db))

	r.POST("/authors", lmsAuthor.CreateAuthor(db))
	r.GET("/authors/:id", lmsAuthor.GetAuthorById(db))
	r.GET("/authors", lmsAuthor.GetAuthorParams(db))
	r.PATCH("/authors/:id", lmsAuthor.UpdateAuthor(db))
	r.DELETE("/authors/:id", lmsAuthor.DeleteAuthor(db))

	r.GET("/genres", lmsGenre.GetGenresByParams(db))
	r.GET("/genres/:id", lmsGenre.GetGenreById(db))
	r.POST("/genres", lmsGenre.CreateGenre(db))
	r.PATCH("/genres/:id", lmsGenre.UpdateGenre(db))
	r.DELETE("genres/:id", lmsGenre.DeleteGenre(db))

	r.POST("/libraries", lmsLibrary.StoreBook(db))
	r.GET("/libraries/:id", lmsLibrary.GetPositionByID(db))
	r.GET("/author/:id", lmsLibrary.GetBooksPositionByAuthor(db))
	r.DELETE("libraries/:id", lmsLibrary.RemoveBook(db))

	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}