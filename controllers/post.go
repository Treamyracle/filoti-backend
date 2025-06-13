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
	// Kita akan tetap mengambil userID karena itu adalah ID admin yang sedang login
	// Tetapi kita TIDAK akan menyimpannya di model Post atau Notification lagi.
	// Ini bisa digunakan untuk 'UpdatedBy' di Status, atau logging internal.
	uidVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var currentUserID uint // Gunakan ini untuk UpdatedBy di Status
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

	// Buat objek Post
	// UserID TIDAK LAGI DITETAPKAN DI SINI karena sudah dihapus dari models.Post
	post := models.Post{
		Title:      input.Title,
		Keterangan: input.Keterangan,
		Ruangan:    input.Ruangan,
		ImageURL:   input.ImageURL,
		ItemType:   input.ItemType,
		// CreatedAt akan diisi otomatis oleh GORM
	}

	// Simpan dalam satu transaksi agar konsisten (post, status, notification)
	tx := config.DB.Begin()
	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post: " + err.Error()})
		return
	}

	// Buat status awal: Status=1 (aktif), UpdatedBy = currentUserID, UpdatedAt otomatis
	status := models.Status{
		PostID:    post.ID,
		Status:    1,             // Default status: 1 (aktif)
		UpdatedBy: currentUserID, // Gunakan ID admin yang sedang login
		// UpdatedAt akan diisi otomatis oleh GORM
	}
	if err := tx.Create(&status).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create status: " + err.Error()})
		return
	}

	// Buat notification: misalnya message dasar
	// UserID TIDAK LAGI DITETAPKAN DI SINI karena sudah dihapus dari models.Notification
	message := "New post created by admin (ID: " + http.StatusText(int(currentUserID)) + "): " + post.Title // Pesan notifikasi
	notification := models.Notification{
		PostID:  post.ID,
		Message: message,
		IsRead:  false, // Default false
		// CreatedAt akan diisi otomatis oleh GORM
	}
	if err := tx.Create(&notification).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification: " + err.Error()})
		return
	}
	tx.Commit()

	// Kembalikan data post (beserta status). Jika ingin menampilkan status:
	// Muat status terkait sebelum respon
	// Preload Status untuk menyertakan data Status dalam respons
	if err := config.DB.Preload("Status").First(&post, post.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post after creation: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully", "post": post})
}

// GetPosts handler: mengambil semua post
func GetPosts(c *gin.Context) {
	var posts []models.Post
	// Preload Status jika ingin sertakan status di response
	// Jika ingin menampilkan username admin, Anda harus:
	// 1. Tambahkan AdminID di models.Post
	// 2. Preload "Admin" di sini.
	if err := config.DB.Preload("Status").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts: " + err.Error()})
		return
	}

	// Jika Anda perlu memodifikasi struktur data sebelum respons (misalnya, menambahkan 'username')
	// Buatlah slice dari map[string]interface{} atau struct kustom untuk respons
	var postsToReturn []gin.H
	for _, p := range posts {
		// Default username, karena UserID tidak lagi ada di model Post
		// Jika Anda menambahkan AdminID ke model Post, Anda bisa preload Admin dan gunakan Admin.Username
		username := "Administrator" // Placeholder

		// Filter berdasarkan itemType dan status di frontend atau di sini jika perlu
		// (Asumsi frontend masih memfilter 'lost' dan 'active')
		// Jika Anda ingin memfilter di backend:
		// if p.ItemType == "lost" && p.Status.Status == 1 { // Assuming 1 means 'active'
		postsToReturn = append(postsToReturn, gin.H{
			"id":         p.ID,
			"username":   username, // Username admin placeholder
			"image_url":  p.ImageURL,
			"title":      p.Title,
			"ruangan":    p.Ruangan,
			"keterangan": p.Keterangan,
			"item_type":  p.ItemType,
			"created_at": p.CreatedAt,
			"status":     p.Status.Status, // Hanya mengembalikan int status
		})
		// }
	}

	// Pastikan selalu mengembalikan array JSON, bahkan jika kosong
	if postsToReturn == nil {
		postsToReturn = []gin.H{}
	}

	c.JSON(http.StatusOK, postsToReturn) // Ubah menjadi array langsung
}

// --- Tambahkan atau sesuaikan handler lain di sini (GetPostByID, UpdatePost, DeletePost, dll.) ---
// ...
