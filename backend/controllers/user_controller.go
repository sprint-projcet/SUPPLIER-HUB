package controllers

import (
	"net/http"
	"strings"

	"supplierhub-backend/config"
	"supplierhub-backend/models"

	"github.com/gin-gonic/gin"
)

func GetUserStats(c *gin.Context) {
	umkmID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var totalOrders int64
	var shippedOrders int64

	// Hitung total pesanan UMKM
	config.DB.Model(&models.Order{}).Where("umkm_id = ?", umkmID).Count(&totalOrders)

	// Hitung pesanan yang sedang dikirim (shipped)
	config.DB.Model(&models.Order{}).Where("umkm_id = ? AND status = ?", umkmID, models.OrderShipped).Count(&shippedOrders)

	c.JSON(http.StatusOK, gin.H{
		"total_orders":   totalOrders,
		"shipped_orders": shippedOrders,
		"vouchers":       0, // Belum ada tabel khusus voucher, set default 0
		"points":         totalOrders * 10, // Contoh logika bisnis: tiap transaksi dapat 10 poin
	})
}

func GetUserOrders(c *gin.Context) {
	umkmID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var orders []models.Order
	// Preload Product agar info detail produk ikut terbawa
	if err := config.DB.Preload("Product").Where("umkm_id = ?", umkmID).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   orders,
	})
}

// --- ALGORITMA INTI ---

// KMPMatch mengimplementasikan Knuth-Morris-Pratt untuk pencarian string
func KMPMatch(text string, pattern string) bool {
	text = strings.ToLower(text)
	pattern = strings.ToLower(pattern)

	n := len(text)
	m := len(pattern)
	if m == 0 {
		return true
	}

	lps := make([]int, m)
	lenLps := 0
	i := 1

	for i < m {
		if pattern[i] == pattern[lenLps] {
			lenLps++
			lps[i] = lenLps
			i++
		} else {
			if lenLps != 0 {
				lenLps = lps[lenLps-1]
			} else {
				lps[i] = 0
				i++
			}
		}
	}

	i = 0
	j := 0
	for i < n {
		if pattern[j] == text[i] {
			j++
			i++
		}
		if j == m {
			return true
		} else if i < n && pattern[j] != text[i] {
			if j != 0 {
				j = lps[j-1]
			} else {
				i++
			}
		}
	}
	return false
}

// QuickSortPrice mengimplementasikan Quick Sort untuk mengurutkan harga
func QuickSortPrice(items []models.Product, low, high int, asc bool) {
	if low < high {
		pi := partitionPrice(items, low, high, asc)
		QuickSortPrice(items, low, pi-1, asc)
		QuickSortPrice(items, pi+1, high, asc)
	}
}

func partitionPrice(items []models.Product, low, high int, asc bool) int {
	pivot := items[high].Price
	i := low - 1
	for j := low; j < high; j++ {
		if asc {
			if items[j].Price <= pivot {
				i++
				items[i], items[j] = items[j], items[i]
			}
		} else {
			if items[j].Price >= pivot {
				i++
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	items[i+1], items[high] = items[high], items[i+1]
	return i + 1
}

// IsItemValid mengimplementasikan Binary Search untuk validasi ketersediaan produk (berdasarkan ID)
func IsItemValid(sortedIDs []string, targetID string) bool {
	low, high := 0, len(sortedIDs)-1
	for low <= high {
		mid := low + (high-low)/2
		if sortedIDs[mid] == targetID {
			return true // Item Valid
		} else if sortedIDs[mid] < targetID {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return false // Item Tidak Ditemukan / Invalid
}

// ----------------------

// GetProducts mengambil katalog produk untuk UMKM dengan fitur Pencarian (KMP) dan Sorting (Quick Sort)
func GetProducts(c *gin.Context) {
	var allProducts []models.Product
	
	// Tarik seluruh data produk dari database
	if err := config.DB.Preload("Supplier").Find(&allProducts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}

	// 1. Filtering menggunakan KMP (jika ada query 'search')
	searchQuery := c.Query("search")
	var filteredProducts []models.Product
	
	if searchQuery != "" {
		for _, product := range allProducts {
			// Mencari kecocokan pattern KMP pada Nama Produk
			if KMPMatch(product.Name, searchQuery) {
				filteredProducts = append(filteredProducts, product)
			}
		}
	} else {
		filteredProducts = allProducts
	}

	// 2. Sorting menggunakan Quick Sort (jika ada query 'sort_by')
	sortBy := c.Query("sort_by")
	if len(filteredProducts) > 0 {
		if sortBy == "price_asc" {
			QuickSortPrice(filteredProducts, 0, len(filteredProducts)-1, true)
		} else if sortBy == "price_desc" {
			QuickSortPrice(filteredProducts, 0, len(filteredProducts)-1, false)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil katalog produk",
		"data":    filteredProducts,
	})
}

type CreateOrderInput struct {
	ItemID   string `json:"item_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

func CreateOrder(c *gin.Context) {
	umkmID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Tarik semua ID produk yang valid dan terurut
	var sortedIDs []string
	config.DB.Model(&models.Product{}).Order("id asc").Pluck("id", &sortedIDs)

	// 2. Validasi ID Produk menggunakan Binary Search (Sesuai Aturan PRD 6.3)
	if !IsItemValid(sortedIDs, input.ItemID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item ID tidak valid atau tidak ditemukan (Binary Search Check Failed)"})
		return
	}

	// 3. Ambil data produk asli dari DB untuk kalkulasi harga
	var product models.Product
	if err := config.DB.First(&product, "id = ?", input.ItemID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}

	// 4. Kalkulasi harga pesanan
	totalBasePrice := product.Price * float64(input.Quantity)
	systemFee := totalBasePrice * 0.03
	grandTotal := totalBasePrice + systemFee

	// 5. Buat dan Simpan Order ke Database
	order := models.Order{
		UmkmID:         umkmID.(string),
		SupplierID:     product.SupplierID,
		ProductID:      product.ID,
		Quantity:       input.Quantity,
		TotalBasePrice: totalBasePrice,
		SystemFee:      systemFee,
		GrandTotal:     grandTotal,
		Status:         models.OrderPending,
	}

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pesanan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Pesanan berhasil dibuat",
		"data":    order,
	})
}
