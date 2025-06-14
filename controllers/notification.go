// controllers/notification_controller.go
package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time" // Untuk formatting waktu jika diperlukan

	"filoti-backend/config"
	"filoti-backend/models"

	"github.com/gin-gonic/gin"
)

// GetNotifications handler: Mengambil notifikasi untuk admin yang sedang login
func GetNotifications(c *gin.Context) {

	var notifications []models.Notification
	// Mengambil semua notifikasi yang terkait dengan Post yang di-upload oleh admin ini,
	// atau semua notifikasi jika notifikasi tidak terkait langsung dengan UserID admin.
	// Karena UserID sudah dihapus dari Notification, kita tidak bisa filter berdasarkan itu.
	// Jika notifikasi ini adalah 'untuk semua admin' atau 'terkait dengan post yang admin buat',
	// Anda perlu preload Post dan filter berdasarkan itu jika perlu.
	// Untuk saat ini, saya akan mengambil semua notifikasi, karena UserID tidak ada di model Notification.
	// Jika notifikasi ditujukan ke admin berdasarkan post yang mereka buat, Anda butuh relasi ke admin di Post.

	// Skenario 1: Notifikasi untuk semua admin (terkait semua post)
	// if err := config.DB.Find(&notifications).Error; err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
	//     return
	// }

	// Skenario 2 (Lebih realistis): Notifikasi yang relevan untuk admin yang login
	// Ini akan membutuhkan `AdminID` di `models.Post` dan relasi di sana.
	// Atau, jika notifikasi adalah untuk setiap post, dan semua admin melihatnya:
	if err := config.DB.Preload("Post").Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications: " + err.Error()})
		return
	}

	// Format notifikasi untuk respons frontend
	var notificationsToReturn []gin.H
	for _, notif := range notifications {
		// Tentukan 'type' notifikasi berdasarkan pesan atau logika lain
		// Ini adalah contoh sederhana, Anda mungkin perlu logika yang lebih canggih
		notificationType := "info"
		iconColor := "bg-blue-500" // Default color
		if notif.IsRead {
			iconColor = "bg-gray-400" // Warna berbeda jika sudah dibaca
		}
		if contains(notif.Message, "baru dibuat") { // Contoh deteksi tipe
			notificationType = "new_post"
			iconColor = "bg-green-500"
		} else if contains(notif.Message, "klaim") || contains(notif.Message, "ambil") {
			notificationType = "claim"
			iconColor = "bg-purple-500"
		} else if contains(notif.Message, "update") || contains(notif.Message, "status") {
			notificationType = "update"
			iconColor = "bg-orange-500"
		}

		notificationsToReturn = append(notificationsToReturn, gin.H{
			"id":         notif.ID,
			"post_id":    notif.PostID,
			"message":    notif.Message, // Ini akan menjadi 'text' di frontend
			"is_read":    notif.IsRead,
			"created_at": notif.CreatedAt,
			"time":       formatTimeAgo(notif.CreatedAt), // Format waktu untuk frontend
			"type":       notificationType,               // Tambahkan 'type' untuk ikon
			"iconColor":  iconColor,                      // Tambahkan 'iconColor'
			"post_title": notif.Post.Title,               // Jika preload "Post" berhasil
		})
	}

	if notificationsToReturn == nil {
		notificationsToReturn = []gin.H{}
	}

	c.JSON(http.StatusOK, notificationsToReturn)
}

// Fungsi helper untuk format waktu (sesuaikan jika perlu)
func formatTimeAgo(t time.Time) string {
	// Implementasi sederhana, bisa lebih canggih
	diff := time.Since(t)
	if diff < time.Minute {
		return "Baru saja"
	} else if diff < time.Hour {
		return fmt.Sprintf("%d menit lalu", int(diff.Minutes()))
	} else if diff < 24*time.Hour {
		return fmt.Sprintf("%d jam lalu", int(diff.Hours()))
	} else if diff < 30*24*time.Hour {
		return fmt.Sprintf("%d hari lalu", int(diff.Hours()/24))
	}
	return t.Format("02 Jan 2006") // Format tanggal jika lebih lama
}

// Fungsi helper sederhana untuk contains string
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// Pastikan Anda import "fmt" dan "strings"
// import (
// 	"fmt"
// 	"strings"
// )
