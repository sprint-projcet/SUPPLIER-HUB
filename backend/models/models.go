package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleUser     Role = "user"
	RoleSupplier Role = "supplier"
	RoleAdmin    Role = "admin"
)

// User merepresentasikan entitas Pengguna, baik itu UMKM, Supplier maupun Admin.
type User struct {
	ID           string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	BusinessName string    `gorm:"type:varchar(255)" json:"business_name"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"` // Disembunyikan di JSON response
	Role         Role      `gorm:"type:varchar(20);not null" json:"role"`
	Address      string    `gorm:"type:text" json:"address"`
	Category     string    `gorm:"type:varchar(100)" json:"category"` // Untuk Kategori Bahan Supplier
	Region       string    `gorm:"type:varchar(100)" json:"region"`   // Untuk Wilayah Supply
	DocumentURL  string    `gorm:"type:varchar(255)" json:"document_url"` // Opsional untuk UMKM, Wajib bagi Supplier
	Status       string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, active, suspended
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	Products []Product `gorm:"foreignKey:SupplierID" json:"products,omitempty"`
	Orders   []Order   `gorm:"foreignKey:UmkmID" json:"purchased_orders,omitempty"`
	Sales    []Order   `gorm:"foreignKey:SupplierID" json:"sales_orders,omitempty"`
}

// Product merepresentasikan barang yang dijual oleh Supplier
type Product struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	SupplierID  string    `gorm:"type:varchar(36);not null;index" json:"supplier_id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Category    string    `gorm:"type:varchar(100)" json:"category"`
	Price       float64   `gorm:"type:numeric(15,2);not null" json:"price"`
	Stock       int       `gorm:"type:int;default:0" json:"stock"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrderStatus defines the lifecycle of an order
type OrderStatus string

const (
	OrderPending    OrderStatus = "pending"
	OrderPaid       OrderStatus = "paid"
	OrderProcessing OrderStatus = "processing"
	OrderShipped    OrderStatus = "shipped"
	OrderCompleted  OrderStatus = "completed"
	OrderCancelled  OrderStatus = "cancelled"
)

// Order merepresentasikan tagihan pembelian antara UMKM dan Supplier
type Order struct {
	ID         string      `gorm:"type:varchar(36);primaryKey" json:"id"`
	UmkmID     string      `gorm:"type:varchar(36);not null;index" json:"umkm_id"`
	SupplierID string      `gorm:"type:varchar(36);not null;index" json:"supplier_id"`
	TotalPrice float64     `gorm:"type:numeric(15,2);not null" json:"total_price"`
	Status     OrderStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// Log merepresentasikan rekam jejak aktivitas (Audit) untuk Admin
type Log struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      string    `gorm:"type:varchar(36);index" json:"user_id"` // Siapa yang melakukan aksi
	Action      string    `gorm:"type:varchar(100);not null" json:"action"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// BeforeCreate hooks untuk generate UUID secara otomatis sebelum insert ke database
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return
}
