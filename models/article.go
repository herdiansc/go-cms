package models

import "github.com/gosimple/slug"

// Article struct
type Article struct {
	Base
	Title                string `gorm:"not null"`
	Content              string `gorm:"not null"`
	Status               string `gorm:"not null"`
	WriterID             int64  `gorm:"not null"`
	Slug                 string `gorm:"not null"`
	TagRelationshipScore int64
}

// CreateArticleRequest struct
type CreateArticleRequest struct {
	Title   string   `json:"title" validate:"required"`
	Content string   `json:"content" validate:"required"`
	Status  string   `json:"status"`
	Tags    []string `json:"tags"`
}

// Article converts CreateArticleRequest to Article
func (c CreateArticleRequest) Article() Article {
	status := "DRAFT"
	if c.Status != "" {
		status = c.Status
	}
	return Article{
		Title:   c.Title,
		Content: c.Content,
		Status:  status,
		Slug:    slug.Make(c.Title),
	}
}

// PatchArticleRequest struct
type PatchArticleRequest struct {
	Status string `json:"status"`
}
