package routes

import (
	"net/http" // Pastikan ini diimpor
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
		sessionSecret = "secret" // sebaiknya override di env
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
		c.Status(http.StatusNoContent) // 204 No Content adalah respons standar untuk OPTIONS
		return
	})
	// -----------------------------------------------------------------------

	// Auth endpoints (ini adalah rute publik, tidak perlu AuthRequired)
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)

	// Grup rute yang memerlukan Autentikasi
	// Semua rute di dalam `authorized` group akan melewati `middleware.AuthRequired()`
	authorized := r.Group("/") // Anda bisa mengganti ini dengan r.Group("/api") jika mau prefix
	authorized.Use(middleware.AuthRequired())
	{
		// Posts routes (sudah ada)
		posts := authorized.Group("/posts") // Pastikan posts ada di bawah grup yang diautentikasi
		{
			posts.POST("", controllers.CreatePost)
			posts.GET("", controllers.GetPosts)
			// ... tambahkan rute posts lainnya seperti GetPostByID, UpdatePost, DeletePost
		}

		// --- TAMBAHKAN RUTE GET /ME DI SINI ---
		authorized.GET("/me", controllers.GetCurrentUser) // Menggunakan GetCurrentUser dari controllers
		// ------------------------------------

		// Anda bisa menambahkan rute terautentikasi lain di sini
		// authorized.GET("/dashboard", controllers.GetDashboardData)
	}

	return r
}
