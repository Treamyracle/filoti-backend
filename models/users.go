package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	IsAdmin   bool      `gorm:"default:false" json:"is_admin"` // Tambahkan field IsAdmin
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
