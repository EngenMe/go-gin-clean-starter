package entity

import (
	"time"

	"gorm.io/gorm"
)

// Timestamp is a struct that provides standard fields for tracking creation, update, and deletion timestamps.
type Timestamp struct {
	CreatedAt time.Time `gorm:"type:timestamp with time zone" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
	DeletedAt gorm.DeletedAt
}

// Authorization represents the structure for user authorization details containing a token and role.
// Token holds the authentication token provided to the user.
// Role defines the access level, restricted to "user" or "admin".
type Authorization struct {
	Token string `json:"token" binding:"required"`
	Role  string `json:"role" binding:"required,oneof=user admin"`
}
