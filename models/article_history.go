package models

// ArticleHistory struct
type ArticleHistory struct {
	Base
	Article   string `gorm:"not null"`
	Version   int64  `gorm:"not null"`
	Status    string `gorm:"not null"`
	ArticleID int64  `gorm:"not null"`
	Action    string `gorm:"not null"`
}
