// models/status.go
package models

import (
	"time"
)

type Status struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	PostID uint `gorm:"uniqueIndex;not null" json:"post_id"`
	Status int  `gorm:"not null;default:1" json:"status"`
	// siapa yang ngambil
	// foto pengambilan
	UpdatedBy uint      `gorm:"not null" json:"updated_by"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	// relasi ke post (optional)
	Post *Post `gorm:"foreignKey:PostID" json:"-"`
}
