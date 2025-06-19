package models

// Role struct
type Role struct {
	Base
	Name  string `gorm:"not null;unique"`
	Title string `gorm:"not null;unique"`
}
