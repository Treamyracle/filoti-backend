// routes/routes.go (Tambahkan di grup yang diautentikasi)

package routes

import (
	"net/http"
	"os"
	"time"

	"filoti-backend/controllers" // Pastikan import controllers sudah ada
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
		AllowOrigins: []string{
			"https://filoti-frontend.vercel.app", // <--- GANTI INI DENGAN DOMAIN FRONTEND VERSEL ANDA!
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

	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
		return
	})

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)

	// Grup rute yang memerlukan Autentikasi
	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired())
	{
		posts := authorized.Group("/posts")
		{
			posts.POST("", controllers.CreatePost)
			posts.GET("", controllers.GetPosts)
			posts.GET("/:id", controllers.GetPostByID)
		}

		authorized.GET("/me", controllers.GetCurrentUser) // Rute untuk mendapatkan user yang login

		// --- TAMBAHKAN RUTE INI ---
		authorized.GET("/notifications", controllers.GetNotifications) // Endpoint baru untuk notifikasi
		// -------------------------
	}

	return r
}
