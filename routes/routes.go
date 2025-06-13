// routes/routes.go
package routes

import (
	"net/http"
	"os"
	"time"

	"filoti-backend/controllers"
	"filoti-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "secret"
	}
	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24,
		HttpOnly: true,
		Secure:   true, // Pastikan ini TRUE di Vercel
		SameSite: http.SameSiteNoneMode,
	})
	r.Use(sessions.Sessions("gin_session", store))

	// CORS middleware
	r.Use(cors.New(cors.Config{
		// PASTIKAN INI MENCAKUP DOMAIN FRONTEND VERSEL ANDA!
		AllowOrigins: []string{
			"https://filoti-frontend.vercel.app", // <--- Contoh domain frontend Vercel Anda
			"http://localhost:5500",
			"http://127.0.0.1:5500",
			"http://localhost:3000",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// --- TAMBAHKAN INI UNTUK MENANGANI PREFLIGHT OPTIONS SECARA EKSPLISIT ---
	// Ini adalah fallback jika middleware CORS tidak menangkapnya dengan benar di Vercel.
	// Ini harus ditempatkan SEBELUM rute spesifik Anda, tapi SETELAH middleware CORS.
	r.OPTIONS("/*path", func(c *gin.Context) {
		// Karena middleware CORS sudah diterapkan, ini hanya perlu merespons OK.
		// Middleware CORS yang akan menambahkan header Access-Control-* yang benar.
		c.Status(http.StatusNoContent) // 204 No Content adalah respons standar untuk OPTIONS
		return
	})
	// -----------------------------------------------------------------------

	// Auth endpoints
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)

	// Posts: but POST requires auth
	posts := r.Group("/posts")
	posts.Use(middleware.AuthRequired())
	{
		posts.POST("", controllers.CreatePost)
		posts.GET("", controllers.GetPosts)
	}

	return r
}
