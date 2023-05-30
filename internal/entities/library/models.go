package library

type Library struct {
	Book_ID   int    `gorm:"not null" json:"book_id"`
	Aisle     int    `json:"aisle"`
	Level     int    `json:"level"`
	Position  int    `json:"position"`
	CreatedAt string `gorm:"create_time" json:"created_at"`
	UpdatedAt string `gorm:"update_time" json:"updated_at"`
	IsDeleted bool   `json:"is_deleted"`
}
