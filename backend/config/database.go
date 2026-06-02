package config

import (
	"database/sql"
	"log"
	"os"
	"time"

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
		err = godotenv.Load("backend/.env")
	}
	if err != nil {
		log.Println("Peringatan: Tidak dapat memuat file .env atau backend/.env (mungkin menggunakan environment default)")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Default config untuk MySQL (contoh development)
		dsn = "root:@tcp(127.0.0.1:3306)/supplierhub?charset=utf8mb4&parseTime=True&loc=Local"
	}

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database MySQL: %v", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Gagal menyiapkan koneksi database MySQL: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := pingDatabase(sqlDB); err != nil {
		log.Fatalf("Gagal memverifikasi koneksi database MySQL: %v. Pastikan MySQL/XAMPP berjalan di 127.0.0.1:3306 dan database supplierhub sudah ada.", err)
	}

	log.Println("Database MySQL Terhubung!")

	// Menjalankan Auto Migration (Menyesuaikan skema tabel ke Data Models secara otomatis)
	// Disabling FK constraint creation here avoids migration failure on legacy
	// data yang memiliki urutan INSERT / orphaned rows sebelumnya.
	err = database.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.Payment{},
		&models.FinanceLog{},
		&models.RequestLog{},
		&models.ShipmentLog{},
		&models.Notification{},
		&models.Wishlist{},
		&models.Log{},
	)
	if err != nil {
		log.Fatalf("Gagal menjalankan migrasi schema database: %v", err)
	}

	DB = database

	// SEEDER: Membuat akun Admin default jika belum ada
	seedAdmin()
}

func pingDatabase(sqlDB *sql.DB) error {
	var err error
	for attempt := 1; attempt <= 5; attempt++ {
		if err = sqlDB.Ping(); err == nil {
			return nil
		}
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	return err
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
