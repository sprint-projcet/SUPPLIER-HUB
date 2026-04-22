package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register menangani pendaftaran awal untuk UMKM dan Supplier
func Register(c *gin.Context) {
	// TODO: Menerima JSON (Email, Password, Role)
	// TODO: Hash password menggunakan bcrypt
	// TODO: Simpan ke PostgreSQL menggunakan GORM
	c.JSON(http.StatusOK, gin.H{
		"message": "Fungsi Register berhasil dipanggil (Stub)",
	})
}

// Login memvalidasi kredensial dan menerbitkan JWT Token
func Login(c *gin.Context) {
	// TODO: Cek email di dalam database
	// TODO: Verifikasi bcrypt password
	// TODO: Jika sukses, hasilkan JWT dengan Claims: id, role
	c.JSON(http.StatusOK, gin.H{
		"message": "Fungsi Login berhasil dipanggil (Stub)",
		"token": "bearer_jwt_token_simulation",
		"role": "user",
	})
}
