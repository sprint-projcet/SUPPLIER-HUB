package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAdminStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"total_suppliers": 128,
		"revenue_growth": "24%",
		"total_transactions": 2400000,
	})
}

func GetAdminSuppliers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Mengembalikan senarai semua supplier"})
}

func VerifySupplier(c *gin.Context) {
	supplierID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Supplier " + supplierID + " berserta perlengkapan statusnya telah diperbarui"})
}

func GetAdminLogs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Log sistem audit keseluruhan dikembalikan"})
}
