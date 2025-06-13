// routes/routes.go
package routes

import (
	"filoti-backend/controllers"
	"filoti-backend/middleware"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Session middleware
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "secret" // sebaiknya override di env
	}
	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteNoneMode, // <--- UBAH KE http.SameSiteNoneMode
	})
	r.Use(sessions.Sessions("gin_session", store))

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

	// Anda bisa menambahkan route lain (misal: GET /me, GET /notifications, dsb.)

	return r
}
