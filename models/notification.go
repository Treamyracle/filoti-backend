package models

import (
	"time"
)

type Notification struct {
	ID uint `gorm:"primaryKey" json:"id"`
	// UserID    uint      `gorm:"not null;index" json:"user_id"` // <-- DIHAPUS sesuai permintaan
	PostID    uint      `gorm:"not null;index" json:"post_id"`
	Message   string    `gorm:"not null" json:"message"`
	IsRead    bool      `gorm:"not null;default:false" json:"is_read"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// relasi opsional
	// User *User `gorm:"foreignKey:UserID" json:"-"` // <-- DIHAPUS sesuai permintaan
	Post *Post `gorm:"foreignKey:PostID" json:"-"`
}
