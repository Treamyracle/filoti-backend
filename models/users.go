// models/user.go
package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	// jika ingin relasi:
	Posts []Post `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}
