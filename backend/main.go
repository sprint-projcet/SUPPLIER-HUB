package main

import (
	"log"

	"supplierhub-backend/config"
	"supplierhub-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inisialisasi Koneksi Database
	config.ConnectDatabase()

	// 2. Setup Gin Router
	r := gin.Default()

	// 3. Konfigurasi Middleware CORS
	// (Mengizinkan UI frontend dari port atau origin berbeda untuk memanggil API ini)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // Boleh diubah ke Origin spesifik UI untuk Production
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// 4. Setup Routes
	routes.SetupRoutes(r)

	// 5. Jalankan Web Server
	port := "8080"
	log.Printf("Server berjalan di port http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
