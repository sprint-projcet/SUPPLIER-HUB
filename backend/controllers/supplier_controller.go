package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"

	"supplierhub-backend/config"
	"supplierhub-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetSupplierStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"stock":      1240,
		"new_orders": 12,
		"revenue_rp": 450000000,
		"rating":     4.8,
	})
}

func GetSupplierProducts(c *gin.Context) {
	supplierID, _ := c.Get("user_id")

	var products []models.Product
	if err := config.DB.Where("supplier_id = ?", supplierID).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func CreateProduct(c *gin.Context) {
	supplierID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	name := c.PostForm("name")
	category := c.PostForm("category")
	priceStr := c.PostForm("price")
	stockStr := c.PostForm("stock")
	description := c.PostForm("description")
	location := c.PostForm("location")

	price, _ := strconv.ParseFloat(priceStr, 64)
	stock, _ := strconv.Atoi(stockStr)

	// File Upload Handling
	file, err := c.FormFile("image")
	var imageURL string
	if err == nil {
		filename := uuid.New().String() + filepath.Ext(file.Filename)
		uploadPath := "uploads/" + filename
		if err := c.SaveUploadedFile(file, uploadPath); err == nil {
			imageURL = config.PublicURL(uploadPath)
		}
	}

	input := models.Product{
		SupplierID:  supplierID.(string),
		Name:        name,
		Category:    category,
		Price:       price,
		Stock:       stock,
		Description: description,
		Location:    location,
		ImageURL:    imageURL,
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Produk berhasil ditambahkan",
		"data":    input,
	})
}

func UpdateProduct(c *gin.Context) {
	supplierID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	productID := c.Param("id")
	var product models.Product
	if err := config.DB.Where("id = ? AND supplier_id = ?", productID, supplierID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan atau Anda tidak berwenang"})
		return
	}

	name := c.PostForm("name")
	category := c.PostForm("category")
	priceStr := c.PostForm("price")
	stockStr := c.PostForm("stock")
	description := c.PostForm("description")
	location := c.PostForm("location")

	if name != "" {
		product.Name = name
	}
	if category != "" {
		product.Category = category
	}
	if priceStr != "" {
		price, _ := strconv.ParseFloat(priceStr, 64)
		product.Price = price
	}
	if stockStr != "" {
		stock, _ := strconv.Atoi(stockStr)
		product.Stock = stock
	}
	if description != "" {
		product.Description = description
	}
	if location != "" {
		product.Location = location
	}

	// File Upload Handling
	file, err := c.FormFile("image")
	if err == nil {
		filename := uuid.New().String() + filepath.Ext(file.Filename)
		uploadPath := "uploads/" + filename
		if err := c.SaveUploadedFile(file, uploadPath); err == nil {
			product.ImageURL = config.PublicURL(uploadPath)
		}
	}

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui produk: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk berhasil diperbarui",
		"data":    product,
	})
}

func DeleteProduct(c *gin.Context) {
	supplierID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	productID := c.Param("id")
	var product models.Product
	if err := config.DB.Where("id = ? AND supplier_id = ?", productID, supplierID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan atau Anda tidak berwenang"})
		return
	}

	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus produk: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk berhasil dihapus",
	})
}

func GetSupplierOrders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Melihat pesanan masuk dari UMKM ke Supplier ini"})
}

func UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Status pesanan " + orderID + " telah diperbarui"})
}
