package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSupplierStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"stock": 1240,
		"new_orders": 12,
		"revenue_rp": 450000000,
		"rating": 4.8,
	})
}

func GetSupplierProducts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Melihat produk milik supplier yang sedang login"})
}

func GetSupplierOrders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Melihat pesanan masuk dari UMKM ke Supplier ini"})
}

func UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Status pesanan " + orderID + " telah diperbarui"})
}
