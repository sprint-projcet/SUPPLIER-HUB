package config

import (
	"log"
	"os"

	"supplierhub-backend/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
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

	// SEEDER: Membuat akun Admin default jika belum ada
	seedAdmin()
}

func seedAdmin() {
	var admin models.User
	// Cek apakah admin sudah ada
	if err := DB.Where("email = ?", "admin@supplierhub.com").First(&admin).Error; err != nil {
		// Jika tidak ada, buat admin baru
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		
		newAdmin := models.User{
			BusinessName: "System Administrator",
			Email:        "admin@supplierhub.com",
			PasswordHash: string(hashedPassword),
			Role:         models.RoleAdmin,
			Status:       "active",
		}
		
		if err := DB.Create(&newAdmin).Error; err == nil {
			log.Println("✅ Akun Admin default berhasil dibuat (admin@supplierhub.com / admin123)")
		} else {
			log.Printf("⚠️ Gagal membuat akun Admin default: %v\n", err)
		}
	}
}
