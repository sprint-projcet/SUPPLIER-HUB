package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"supplierhub-backend/config"
	"supplierhub-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// DTO untuk input form Register
type RegisterInput struct {
	BusinessName string `json:"business_name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	Role         string `json:"role" binding:"required"`
}

// DTO untuk input form Login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register menangani pendaftaran awal untuk UMKM dan Supplier
func Register(c *gin.Context) {
	// Parse fields dari form-data
	businessName := c.PostForm("business_name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	role := c.PostForm("role")
	address := c.PostForm("address")
	category := c.PostForm("category")
	region := c.PostForm("region")

	// 1. Validasi Input
	if businessName == "" || email == "" || password == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kolom Nama, Email, Password, dan Role wajib diisi!"})
		return
	}

	// 2. Cek apakah email sudah terdaftar
	var existingUser models.User
	if err := config.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email sudah terdaftar!"})
		return
	}

	// 3. File Upload Handling (Dokumen)
	var documentURL string
	file, err := c.FormFile("document")
	if err == nil {
		// Buat folder jika belum ada
		os.MkdirAll("uploads/documents", os.ModePerm)
		
		// Gunakan timestamp untuk nama file yang unik
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		filepath := "uploads/documents/" + filename
		
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan dokumen legalitas"})
			return
		}
		documentURL = filepath
	} else if role == "supplier" {
		// Supplier wajib upload dokumen
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dokumen legalitas (SIUP/Akta) wajib diunggah untuk pendaftaran supplier"})
		return
	}

	// 4. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}

	// 5. Simpan ke Database
	newUser := models.User{
		BusinessName: businessName,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         models.Role(role),
		Address:      address,
		Category:     category,
		Region:       region,
		DocumentURL:  documentURL,
		Status:       "pending", // Default status, mungkin butuh verifikasi admin nanti
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data pengguna"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil!",
		"user": gin.H{
			"id":            newUser.ID,
			"business_name": newUser.BusinessName,
			"email":         newUser.Email,
			"role":          newUser.Role,
		},
	})
}

// Login memvalidasi kredensial dan menerbitkan JWT Token
func Login(c *gin.Context) {
	var input LoginInput

	// 1. Validasi Input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid atau kurang lengkap"})
		return
	}

	// 2. Cari user berdasarkan Email
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// 3. Verifikasi Password dengan Hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// 4. Cek Status User (Optional: Jangan ijinkan login jika status suspended)
	if user.Status == "suspended" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun Anda ditangguhkan"})
		return
	}

	// 5. Generate JWT Token
	// Ambil secret dari ENV, atau gunakan default secret jika ENV tidak ada (hanya untuk development)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "super_secret_key_supplierhub"
	}

	// Buat claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expired dalam 72 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token autentikasi"})
		return
	}

	// 6. Kirim Respons Berhasil
	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString,
		"role":    user.Role,
		"user": gin.H{
			"id":            user.ID,
			"business_name": user.BusinessName,
			"email":         user.Email,
		},
	})
}
