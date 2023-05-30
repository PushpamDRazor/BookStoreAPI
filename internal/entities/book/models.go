package book

type Book struct {
	ID              int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string `gorm:"not null" json:"title"`
	AuthorName      string `gorm:"not null" json:"author_name"`
	AuthorId        int    `gorm:"not null" json:"author_id"`
	GenreName       string `gorm:"not null" json:"genre_name"`
	GenreId         int    `gorm:"not null" json:"genre_id"`
	PublicationDate string `gorm:"not null" json:"publication_date"`
	CreatedAt       string `gorm:"create_time" json:"created_at"`
	UpdatedAt       string `gorm:"update_time" json:"updated_at"`
	IsDeleted       bool   `json:"is_deleted"`
}
