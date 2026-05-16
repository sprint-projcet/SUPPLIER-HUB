package controllers

import (
	"net/http"
	"sort"
	"strings"

	"supplierhub-backend/config"
	"supplierhub-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getAuthenticatedUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return "", false
	}

	userIDString, ok := userID.(string)
	if !ok || userIDString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi pengguna tidak valid"})
		return "", false
	}

	return userIDString, true
}

func GetUserStats(c *gin.Context) {
	umkmID, ok := getAuthenticatedUserID(c)
	if !ok {
		return
	}

	var totalOrders int64
	var pendingOrders int64
	var shippedOrders int64
	var completedOrders int64
	var totalSpending float64

	// Hitung total pesanan UMKM
	config.DB.Model(&models.Order{}).Where("umkm_id = ?", umkmID).Count(&totalOrders)
	config.DB.Model(&models.Order{}).Where("umkm_id = ? AND status = ?", umkmID, models.OrderPending).Count(&pendingOrders)

	// Hitung pesanan yang sedang dikirim (shipped)
	config.DB.Model(&models.Order{}).Where("umkm_id = ? AND status = ?", umkmID, models.OrderShipped).Count(&shippedOrders)
	config.DB.Model(&models.Order{}).Where("umkm_id = ? AND status = ?", umkmID, models.OrderCompleted).Count(&completedOrders)
	config.DB.Model(&models.Order{}).
		Where("umkm_id = ? AND status IN ?", umkmID, []models.OrderStatus{models.OrderPaid, models.OrderProcessing, models.OrderShipped, models.OrderCompleted}).
		Select("COALESCE(SUM(grand_total), 0)").
		Scan(&totalSpending)

	c.JSON(http.StatusOK, gin.H{
		"total_orders":     totalOrders,
		"pending_orders":   pendingOrders,
		"shipped_orders":   shippedOrders,
		"completed_orders": completedOrders,
		"total_spending":   totalSpending,
		"vouchers":         0,                // Belum ada tabel khusus voucher, set default 0
		"points":           totalOrders * 10, // Contoh logika bisnis: tiap transaksi dapat 10 poin
	})
}

func GetUserOrders(c *gin.Context) {
	umkmID, ok := getAuthenticatedUserID(c)
	if !ok {
		return
	}

	var orders []models.Order
	query := config.DB.Preload("Product").Preload("Product.Supplier").Where("umkm_id = ?", umkmID)

	if status := c.Query("status"); status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	// Preload Product agar info detail produk ikut terbawa
	if err := query.Order("created_at DESC").Find(&orders).Error; err != nil {
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
			// Mencari kecocokan pattern KMP pada nama produk, kategori, lokasi, atau nama supplier
			if KMPMatch(product.Name, searchQuery) ||
				KMPMatch(product.Category, searchQuery) ||
				KMPMatch(product.Location, searchQuery) ||
				KMPMatch(product.Supplier.BusinessName, searchQuery) {
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
	umkmID, ok := getAuthenticatedUserID(c)
	if !ok {
		return
	}

	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Tarik semua ID produk yang valid dan pastikan diurutkan dengan aturan Go (ASCII)
	var sortedIDs []string
	if err := config.DB.Model(&models.Product{}).Pluck("id", &sortedIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memvalidasi produk"})
		return
	}
	sort.Strings(sortedIDs)

	// 2. Validasi ID Produk menggunakan Binary Search (Sesuai Aturan PRD 6.3)
	if !IsItemValid(sortedIDs, input.ItemID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item ID tidak valid atau tidak ditemukan (Binary Search Check Failed)"})
		return
	}

	// 3. Ambil data produk asli dari DB untuk kalkulasi harga
	var product models.Product
	if err := config.DB.First(&product, "id = ?", input.ItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	if product.Stock < input.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stok produk tidak mencukupi"})
		return
	}

	// 4. Kalkulasi harga pesanan
	totalBasePrice := product.Price * float64(input.Quantity)
	systemFee := totalBasePrice * 0.03
	grandTotal := totalBasePrice + systemFee

	// 5. Buat dan Simpan Order ke Database
	order := models.Order{
		UmkmID:         umkmID,
		SupplierID:     product.SupplierID,
		ProductID:      product.ID,
		Quantity:       input.Quantity,
		TotalBasePrice: totalBasePrice,
		SystemFee:      systemFee,
		GrandTotal:     grandTotal,
		Status:         models.OrderPending,
	}

	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		return tx.Model(&product).Update("stock", product.Stock-input.Quantity).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pesanan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Pesanan berhasil dibuat",
		"data":    order,
	})
}

func CancelOrder(c *gin.Context) {
	umkmID, ok := getAuthenticatedUserID(c)
	if !ok {
		return
	}

	orderID := c.Param("id")
	var order models.Order
	if err := config.DB.Where("id = ? AND umkm_id = ?", orderID, umkmID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	if order.Status != models.OrderPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pesanan hanya bisa dibatalkan saat status pending"})
		return
	}

	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Product{}).
			Where("id = ?", order.ProductID).
			Update("stock", gorm.Expr("stock + ?", order.Quantity)).Error; err != nil {
			return err
		}
		return tx.Model(&order).Update("status", models.OrderCancelled).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membatalkan pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pesanan berhasil dibatalkan",
	})
}

func CompleteOrder(c *gin.Context) {
	umkmID, ok := getAuthenticatedUserID(c)
	if !ok {
		return
	}

	orderID := c.Param("id")
	var order models.Order
	if err := config.DB.Where("id = ? AND umkm_id = ?", orderID, umkmID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	if order.Status != models.OrderShipped && order.Status != models.OrderPaid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pesanan belum dapat dikonfirmasi selesai"})
		return
	}

	if err := config.DB.Model(&order).Update("status", models.OrderCompleted).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyelesaikan pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pesanan berhasil dikonfirmasi selesai",
	})
}
