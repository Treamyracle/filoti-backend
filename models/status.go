// models/status.go
package models

import (
	"time"
)

type Status struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PostID      uint      `gorm:"uniqueIndex;not null" json:"post_id"` // Pastikan PostID unik untuk status
	Status      int       `gorm:"default:1" json:"status"`             // 1: Active, 0: Done (atau angka lain sesuai kebutuhan)
	ClaimerName string    `json:"claimer_name,omitempty"`              // Nama pengambil/penemu, opsional
	ProofImage  string    `json:"proof_image,omitempty"`               // URL bukti gambar, opsional
	UpdatedBy   uint      `json:"updated_by"`                          // ID user yang mengubah status
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	// relasi ke post (optional)
	Post *Post `gorm:"foreignKey:PostID" json:"-"`
}
