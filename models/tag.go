package models

import "strings"

// Tag struct
type Tag struct {
	Base
	Title string `gorm:"not null;unique"`
}

// CreateTagRequest struct
type CreateTagRequest struct {
	Title string `json:"title" validate:"required"`
}

// Tag converts CreateTagRequest to Tag
func (c CreateTagRequest) Tag() Tag {
	return Tag{
		Title: strings.ToLower(c.Title),
	}
}

// TagUsageQueryResult struct
type TagUsageQueryResult struct {
	ID    int64
	Count int64
}

// TagListItem struct
type TagListItem struct {
	ID         int64  `json:"id"`
	UUID       string `json:"uuid"`
	Title      string `json:"title"`
	UsageCount int64  `json:"usageCount"`
}

// TagDetail struct
type TagDetail struct {
	Tag        Tag
	UsageCount int64 `json:"usageCount"`
}
