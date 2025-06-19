package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PublicBase struct
type PublicBase struct {
	UUID      string     `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// Base struct
type Base struct {
	ID int64 `gorm:"autoIncrement"`
	PublicBase
}

// BeforeCreate creates a UUID.
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.UUID = uuid.New().String()
	return
}

// Response struct
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type authVerifyCtxType string

const AuthVerifyCtxKey authVerifyCtxType = "authVerifyCtxKey"
