package author

type Author struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	AuthorName  string `gorm:"not null" json:"author_name"`
	Biography   string `json:"biography"`
	DateOfBirth string `json:"date_of_birth"`
	Nationality string `json:"nationality"`
	CreatedAt   string `gorm:"create_time" json:"created_at"`
	UpdatedAt   string `gorm:"update_time" json:"updated_at"`
	IsDeleted   bool   `json:"is_deleted"`
}
