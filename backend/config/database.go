package config

import (
	"log"
	"os"

	"supplierhub-backend/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase menginisialisasi hubungan koneksi aplikasi dengan MySQL
func ConnectDatabase() {
	// Memuat konfigurasi environment variables (opsional jika sudah ada OS ENV)
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Tidak dapat memuat file .env (mungkin menggunakan environment default)")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Default config untuk MySQL (contoh development)
		dsn = "root:@tcp(127.0.0.1:3306)/supplierhub?charset=utf8mb4&parseTime=True&loc=Local"
	}

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database MySQL: %v", err)
	}

	log.Println("Database MySQL Terhubung!")

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
