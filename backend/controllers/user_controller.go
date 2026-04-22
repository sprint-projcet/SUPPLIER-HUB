package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"total_orders": 24,
		"shipped_orders": 3,
		"vouchers": 5,
		"points": 12500,
	})
}

func GetUserOrders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Mengembalikan riwayat pembelian pengguna"})
}

func GetProducts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Melihat katalog produk aktif dari supplier"})
}

func CreateOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Membuat pesanan baru"})
}
