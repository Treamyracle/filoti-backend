// controllers/post.go
package controllers

import (
	"net/http"

	"filoti-backend/config"
	"filoti-backend/models"

	"github.com/gin-gonic/gin"
)

// Input untuk membuat post
type CreatePostInput struct {
	Title      string `json:"title" binding:"required"`
	Keterangan string `json:"keterangan" binding:"required"`
	Ruangan    string `json:"ruangan" binding:"required"`
	ImageURL   string `json:"image_url" binding:"required"`
	ItemType   string `json:"itemType" binding:"required"`
}

// CreatePost handler: memerlukan AuthRequired middleware
func CreatePost(c *gin.Context) {
	// Ambil userID dari context (diset middleware AuthRequired)
	uidVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// ubah ke uint
	var userID uint
	switch v := uidVal.(type) {
	case uint:
		userID = v
	case int:
		userID = uint(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		return
	}

	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Buat objek Post
	post := models.Post{
		UserID:     userID,
		Title:      input.Title,
		Keterangan: input.Keterangan,
		Ruangan:    input.Ruangan,
		ImageURL:   input.ImageURL,
		ItemType:   input.ItemType,
	}
	// Simpan dalam satu transaksi agar konsisten (post, status, notification)
	tx := config.DB.Begin()
	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}
	// Buat status awal: Status=1 (aktif), UpdatedBy = userID, UpdatedAt otomatis
	status := models.Status{
		PostID:    post.ID,
		Status:    1,
		UpdatedBy: userID,
	}
	if err := tx.Create(&status).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create status"})
		return
	}
	// Buat notification: misalnya message dasar
	message := "New post created: " + post.Title
	notification := models.Notification{
		UserID:  userID,
		PostID:  post.ID,
		Message: message,
		IsRead:  false,
	}
	if err := tx.Create(&notification).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
		return
	}
	tx.Commit()

	// Kembalikan data post (beserta status). Jika ingin menampilkan status:
	// Muat status terkait sebelum respon
	config.DB.Preload("Status").First(&post, post.ID)
	c.JSON(http.StatusCreated, gin.H{"post": post})
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	// Preload Status jika ingin sertakan status di response
	if err := config.DB.Preload("Status").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}
