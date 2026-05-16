package controllers

import (
	"errors"
	"net/http"
	"strings"

	"supplierhub-backend/config"
	"supplierhub-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type adminSupplierResponse struct {
	ID           string `json:"id"`
	BusinessName string `json:"business_name"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	Category     string `json:"category"`
	Region       string `json:"region"`
	DocumentURL  string `json:"document_url"`
	Status       string `json:"status"`
	ProductCount int64  `json:"product_count"`
}

func GetAdminStats(c *gin.Context) {
	var totalSuppliers int64
	var totalTransactions int64
	var activeOrders int64
	var pendingSuppliers int64
	var totalRevenue float64
	var systemFeeRevenue float64
	var recentOrders []models.Order

	activeStatuses := []models.OrderStatus{
		models.OrderPending,
		models.OrderPaid,
		models.OrderProcessing,
		models.OrderShipped,
	}

	if err := config.DB.Model(&models.User{}).Where("role = ?", models.RoleSupplier).Count(&totalSuppliers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung total supplier"})
		return
	}
	if err := config.DB.Model(&models.User{}).Where("role = ? AND status = ?", models.RoleSupplier, "pending").Count(&pendingSuppliers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung supplier pending"})
		return
	}
	if err := config.DB.Model(&models.Order{}).Count(&totalTransactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung total transaksi"})
		return
	}
	if err := config.DB.Model(&models.Order{}).Where("status IN ?", activeStatuses).Count(&activeOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung pesanan aktif"})
		return
	}
	if err := config.DB.Model(&models.Order{}).
		Where("status <> ?", models.OrderCancelled).
		Select("COALESCE(SUM(grand_total), 0)").
		Scan(&totalRevenue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung total revenue"})
		return
	}
	if err := config.DB.Model(&models.Order{}).
		Where("status <> ?", models.OrderCancelled).
		Select("COALESCE(SUM(system_fee), 0)").
		Scan(&systemFeeRevenue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung biaya layanan"})
		return
	}
	if err := config.DB.Preload("Product").Order("created_at DESC").Limit(5).Find(&recentOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil transaksi terbaru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":             "success",
		"total_suppliers":    totalSuppliers,
		"pending_suppliers":  pendingSuppliers,
		"total_transactions": totalTransactions,
		"active_orders":      activeOrders,
		"total_revenue":      totalRevenue,
		"system_fee_revenue": systemFeeRevenue,
		"revenue_growth":     "0%",
		"recent_orders":      recentOrders,
	})
}

func GetAdminSuppliers(c *gin.Context) {
	var suppliers []models.User
	query := config.DB.Where("role = ?", models.RoleSupplier)

	if status := strings.TrimSpace(c.Query("status")); status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	if search := strings.TrimSpace(c.Query("search")); search != "" {
		likeSearch := "%" + search + "%"
		query = query.Where(
			"business_name LIKE ? OR email LIKE ? OR category LIKE ? OR region LIKE ?",
			likeSearch,
			likeSearch,
			likeSearch,
			likeSearch,
		)
	}

	if err := query.Order("created_at DESC").Find(&suppliers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar supplier"})
		return
	}

	response := make([]adminSupplierResponse, 0, len(suppliers))
	for _, supplier := range suppliers {
		var productCount int64
		if err := config.DB.Model(&models.Product{}).Where("supplier_id = ?", supplier.ID).Count(&productCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung produk supplier"})
			return
		}

		response = append(response, adminSupplierResponse{
			ID:           supplier.ID,
			BusinessName: supplier.BusinessName,
			Email:        supplier.Email,
			Address:      supplier.Address,
			Category:     supplier.Category,
			Region:       supplier.Region,
			DocumentURL:  supplier.DocumentURL,
			Status:       supplier.Status,
			ProductCount: productCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}

func VerifySupplier(c *gin.Context) {
	supplierID := c.Param("id")
	var supplier models.User

	if err := config.DB.Where("id = ? AND role = ?", supplierID, models.RoleSupplier).First(&supplier).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier tidak ditemukan"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data supplier"})
		return
	}

	if err := config.DB.Model(&supplier).Update("status", "active").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi supplier"})
		return
	}

	adminID, _ := c.Get("user_id")
	_ = config.DB.Create(&models.Log{
		UserID:      toString(adminID),
		Action:      "VERIFY_SUPPLIER",
		Description: "Supplier " + supplier.BusinessName + " diverifikasi oleh admin",
	}).Error

	supplier.Status = "active"
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Supplier berhasil diverifikasi",
		"data": gin.H{
			"id":            supplier.ID,
			"business_name": supplier.BusinessName,
			"email":         supplier.Email,
			"status":        supplier.Status,
		},
	})
}

func GetAdminLogs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Log sistem audit keseluruhan dikembalikan"})
}

func toString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}

	return ""
}
