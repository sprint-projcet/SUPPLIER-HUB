package config

import (
	"log"
	"os"

	"supplierhub-backend/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase menginisialisasi hubungan koneksi aplikasi dengan PostgreSQL
func ConnectDatabase() {
	// Memuat konfigurasi environment variables (opsional jika sudah ada OS ENV)
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Tidak dapat memuat file .env (mungkin menggunakan environment default)")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Default config (contoh development)
		dsn = "host=localhost user=postgres password=postgres dbname=supplierhub port=5432 sslmode=disable"
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	log.Println("Database Terhubung!")

	// Menjalankan Auto Migration (Menyesuaikan skema tabel ke Data Models secara otomatis)
	err = database.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.Log{},
	)
	if err != nil {
		log.Fatalf("Gagal menjalankan migrasi schema database: %v", err)
	}

	DB = database
}
