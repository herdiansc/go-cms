package models

import (
	"time"
)

// PublicBase struct
type PublicBase struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// Base struct
type Base struct {
	ID int64 `gorm:"autoIncrement"`
	PublicBase
}

// Response struct
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type authVerifyCtxType string

const AuthVerifyCtxKey authVerifyCtxType = "authVerifyCtxKey"
