package controllers

import (
	"fmt" // <--- PASTIKAN INI DIIMPOR untuk fmt.Sprintf
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
	uidVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var currentUserID uint
	switch v := uidVal.(type) {
	case uint:
		currentUserID = v
	case int: // Untuk jaga-jaga jika GORM menyimpannya sebagai int
		currentUserID = uint(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		return
	}

	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := models.Post{
		Title:      input.Title,
		Keterangan: input.Keterangan,
		Ruangan:    input.Ruangan,
		ImageURL:   input.ImageURL,
		ItemType:   input.ItemType,
	}

	tx := config.DB.Begin()
	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post: " + err.Error()})
		return
	}

	status := models.Status{
		PostID:    post.ID,
		Status:    1,
		UpdatedBy: currentUserID, // Gunakan ID admin yang sedang login
	}
	if err := tx.Create(&status).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create status: " + err.Error()})
		return
	}

	// PERBAIKAN DI SINI: Gunakan fmt.Sprintf untuk memformat ID menjadi string
	message := fmt.Sprintf("Post baru dibuat oleh admin (ID: %d): %s", currentUserID, post.Title) // <--- PERBAIKAN PENTING

	notification := models.Notification{
		PostID:  post.ID,
		Message: message,
		IsRead:  false,
	}
	if err := tx.Create(&notification).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification: " + err.Error()})
		return
	}
	tx.Commit()

	if err := config.DB.Preload("Status").First(&post, post.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post after creation: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully", "post": post})
}

// GetPosts handler: mengambil semua post
func GetPosts(c *gin.Context) {
	var posts []models.Post
	if err := config.DB.Preload("Status").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts: " + err.Error()})
		return
	}

	var postsToReturn []gin.H
	for _, p := range posts {
		username := "Administrator" // Placeholder

		postsToReturn = append(postsToReturn, gin.H{
			"id":         p.ID,
			"username":   username,
			"image_url":  p.ImageURL,
			"title":      p.Title,
			"ruangan":    p.Ruangan,
			"keterangan": p.Keterangan,
			"item_type":  p.ItemType,
			"created_at": p.CreatedAt,
			"status":     p.Status.Status,
		})
	}

	if postsToReturn == nil {
		postsToReturn = []gin.H{}
	}

	c.JSON(http.StatusOK, postsToReturn)
}
