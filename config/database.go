package config

import (
	"fmt"
	"log"
	"os"

	"filoti-backend/models" // Pastikan path ini benar

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	// Memuat .env jika ada
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Ambil detail koneksi dari environment variables untuk pooler
	// CATATAN: user di Supabase pooler seringkali menyertakan project ref,
	// contoh: "postgres.ugcabvgyvarjkwxifhxz"
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST_POOLER") // Menggunakan nama variabel baru untuk host pooler
	port := os.Getenv("DB_PORT_POOLER") // Menggunakan nama variabel baru untuk port pooler
	dbname := os.Getenv("DB_NAME")

	// Pastikan semua variabel lingkungan penting sudah diatur
	if user == "" || password == "" || host == "" || port == "" || dbname == "" {
		log.Fatal("One or more required database environment variables (DB_USER, DB_PASSWORD, DB_HOST_POOLER, DB_PORT_POOLER, DB_NAME) are not set. Please check your .env file.")
	}

	// Sslmode untuk pooler biasanya "disable" atau "require" tergantung setup Anda.
	// Supabase seringkali merekomendasikan "prefer" atau "require" untuk direct,
	// tapi pooler mungkin berbeda. Coba "disable" dulu jika "require" bermasalah.
	// Jika ini untuk development lokal, "disable" seringkali lebih mudah.
	sslmode := os.Getenv("DB_SSLMODE_POOLER") // Variabel baru untuk sslmode pooler
	if sslmode == "" {
		sslmode = "disable" // Default untuk pooler di development. Ganti ke "require" di produksi jika aman.
	}

	// Buat DSN untuk Connection Pooler
	// Penting: Sertakan parameter `pgbouncer=true` di akhir DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta&pgbouncer=true",
		host, user, password, dbname, port, sslmode,
	)

	log.Printf("Attempting to connect to database using DSN (Pooler): %s", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log semua kueri dan operasi GORM
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db

	// Opsional: Ping database untuk memastikan koneksi aktif
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database!")

	// AutoMigrate models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Status{},
		&models.Notification{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	log.Println("Database connected and migrated successfully.")
}
