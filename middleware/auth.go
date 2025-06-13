// middleware/auth.go
package middleware

import (
	"net/http"

	"filoti-backend/config"
	"filoti-backend/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid := session.Get("id")
		if uid == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		// Cast to uint
		userID, ok := uid.(uint)
		if !ok {
			// mungkin disimpan sebagai int?
			if tmpInt, ok2 := uid.(int); ok2 {
				userID = uint(tmpInt)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
				c.Abort()
				return
			}
		}
		// Opsional: muat user dari DB untuk memastikan masih ada
		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		// Simpan userID (dan user) di context
		c.Set("userID", userID)
		c.Next()
	}
}
