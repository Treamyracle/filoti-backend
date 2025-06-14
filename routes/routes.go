package routes

import (
	"net/http"
	"os"
	"time"

	"filoti-backend/controllers"
	"filoti-backend/middleware" // Menggunakan 'middleware' sesuai dengan kode Anda

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
		Secure:   true, // Pastikan ini TRUE di Vercel jika menggunakan HTTPS
		SameSite: http.SameSiteNoneMode,
	})
	r.Use(sessions.Sessions("gin_session", store))

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
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

	// Handle OPTIONS preflight requests
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// --- Public Routes (Dapat diakses tanpa autentikasi) ---
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/guest-login", controllers.GuestLogin)
	r.GET("/locations", controllers.GetUniqueLocations)
	r.GET("/posts", controllers.GetPosts)        // Postingan dapat dilihat oleh siapa saja
	r.GET("/posts/:id", controllers.GetPostByID) // Detail postingan dapat dilihat oleh siapa saja

	// --- Authenticated Routes (Memerlukan session yang valid) ---
	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired()) // Menggunakan middleware dari paket 'middleware'
	{
		// Rute untuk user yang sedang login
		authorized.GET("/me", controllers.GetCurrentUser)
		authorized.POST("/logout", controllers.Logout)
		authorized.GET("/notifications", controllers.GetNotifications) // Endpoint notifikasi

		// Rute Post yang memerlukan autentikasi (dan cek isAdmin di controllernya)
		posts := authorized.Group("/posts")
		{
			posts.POST("", controllers.CreatePost)             // Membuat post (memerlukan login)
			posts.PUT("/:id", controllers.UpdatePost)          // Update post (memerlukan login & isAdmin)
			posts.DELETE("/:id", controllers.DeletePost)       // Hapus post (memerlukan login & isAdmin)
			posts.PUT("/:id/done", controllers.MarkPostAsDone) // Tandai selesai (memerlukan login & isAdmin)
		}
	}

	return r
}
