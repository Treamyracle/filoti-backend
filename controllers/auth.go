package controllers

import (
	"log" // Tambahkan import log
	"net/http"
	"strings" // Tambahkan import strings

	"filoti-backend/config"
	"filoti-backend/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignupInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Signup(c *gin.Context) {
	var input SignupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Signup: Bad request - %v", err) // Log error binding
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalisasi username
	input.Username = strings.ToLower(strings.TrimSpace(input.Username))
	input.Password = strings.TrimSpace(input.Password) // Pangkas spasi password juga

	log.Printf("Signup: Attempting to signup user: %s", input.Username)

	// Periksa username unik
	var existing models.User
	if err := config.DB.Where("username = ?", input.Username).First(&existing).Error; err == nil {
		log.Printf("Signup: Username '%s' already taken.", input.Username)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
		return
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Signup: Failed to hash password for %s - %v", input.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	log.Printf("Signup: Password hashed for %s", input.Username)

	user := models.User{
		Username: input.Username,
		Password: string(hashed),
	}
	if err := config.DB.Create(&user).Error; err != nil {
		log.Printf("Signup: Failed to create user %s in DB - %v", input.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	log.Printf("Signup: User '%s' created successfully with ID %d.", user.Username, user.ID)

	c.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"created_at": user.CreatedAt,
	})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Login: Bad request - %v", err) // Log error binding
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalisasi username
	input.Username = strings.ToLower(strings.TrimSpace(input.Username))
	input.Password = strings.TrimSpace(input.Password) // Pangkas spasi password juga

	log.Printf("Login: Attempting to login user: %s", input.Username)

	var user models.User
	// Tambahkan logging untuk kueri Find
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		// Jika user tidak ditemukan, GORM akan mengembalikan gorm.ErrRecordNotFound
		log.Printf("Login: User '%s' not found or DB error - %v", input.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	log.Printf("Login: User '%s' found. Comparing passwords...", user.Username)

	// Compare password
	// Log hash yang diambil dari DB dan password yang diinput
	log.Printf("Login: Hashed from DB: '%s', Input password: '%s'", user.Password, input.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		log.Printf("Login: Password mismatch for user '%s' - %v", user.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	log.Printf("Login: Password matched for user '%s'.", user.Username)

	// Set session cookie
	session := sessions.Default(c)
	session.Set("id", int(user.ID))
	if err := session.Save(); err != nil {
		log.Printf("Login: Failed to save session for user %s - %v", user.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	log.Printf("Login: User '%s' logged in successfully. Session ID: %v", user.Username, user.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1}) // expire immediately
	session.Save()
	log.Println("User logged out.")
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
