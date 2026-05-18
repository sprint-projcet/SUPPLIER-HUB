package controllers

import (
	"encoding/json"
	"net/http"

	"supplierhub-backend/config"
	"supplierhub-backend/models"

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

// GetPublicCatalog mengambil semua produk beserta data supplier untuk ditampilkan di halaman katalog publik
func GetPublicCatalog(c *gin.Context) {
	var products []models.Product
	if err := config.DB.Preload("Supplier").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data katalog produk"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func CreateOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Anda harus login"})
		return
	}

	var cart []struct {
		ID         string  `json:"id"`
		SupplierID string  `json:"supplier_id"`
		Title      string  `json:"title"`
		Price      float64 `json:"price"`
		Qty        int     `json:"qty"`
	}

	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data keranjang tidak valid"})
		return
	}

	if len(cart) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Keranjang belanja kosong"})
		return
	}

	// Kelompokkan pesanan berdasarkan Supplier
	ordersBySupplier := make(map[string][]interface{})
	totalsBySupplier := make(map[string]float64)

	for _, item := range cart {
		if item.SupplierID == "" {
			continue // Skip produk tanpa supplier valid
		}
		ordersBySupplier[item.SupplierID] = append(ordersBySupplier[item.SupplierID], item)
		totalsBySupplier[item.SupplierID] += item.Price * float64(item.Qty)
	}

	// Buat entitas Order di database untuk tiap Supplier yang terlibat
	for suppID, items := range ordersBySupplier {
		itemsJSON, _ := json.Marshal(items)

		order := models.Order{
			UmkmID:     userID.(string),
			SupplierID: suppID,
			TotalPrice: totalsBySupplier[suppID],
			OrderItems: string(itemsJSON),
			Status:     models.OrderPending,
		}

		if err := config.DB.Create(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Sebagian pesanan gagal diproses: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Semua pesanan berhasil dibuat dan diteruskan ke Supplier terkait!"})
}
