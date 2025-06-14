package models

import (
	"time"
)

type Post struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ImageURL   string    `json:"image_url"`
	Title      string    `json:"title"`
	Ruangan    string    `json:"ruangan"`
	Keterangan string    `json:"keterangan"`
	ItemType   string    `json:"itemType"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`

	Status        Status         `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"status"` // Status sebagai relasi satu-ke-satu
	Notifications []Notification `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
}
