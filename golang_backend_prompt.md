# Prompt Pengembangan Backend Golang & Database untuk SUPPLIER-HUB

Gunakan prompt di bawah ini kepada AI / developer backend untuk menginstruksikan pembuatan sistem backend Golang yang memfasilitasi aplikasi frontend SUPPLIER-HUB:

---

## 💻 PROMPT UNTUK BACKEND ENGINEER / AI

**Tujuan:**
Buatkan backend berbasis **Golang** untuk aplikasi B2B e-commerce bernama "SUPPLIER-HUB". Aplikasi ini menghubungkan 3 jenis role pengguna: `Umkm` (Pembeli/User), `Supplier` (Penjual), dan `Admin` (Superuser). Backend harus mendukung autentikasi berbasis JWT, manajemen produk, manajemen pesanan, dan analitik data untuk dashboard masing-masing role.

### 1. Stack Teknologi yang Diinginkan
- **Bahasa:** Golang (Go 1.21+)
- **Framework Web:** Gin atau Fiber (pilih salah satu yang paling optimal)
- **Database:** PostgreSQL
- **ORM:** GORM
- **Autentikasi:** JWT (JSON Web Token) dengan pengelolaan Cookie/Header Authorization
- **Lain-lain:** Bcrypt untuk hashing password, environment variables (godotenv).

### 2. Arsitektur Database (Entity-Relationship)
Buatkan model GORM untuk tabel-tabel berikut:

- **Users:**
  - `id` (UUID, Primary Key)
  - `business_name` (String)
  - `email` (String, Unique)
  - `password_hash` (String)
  - `role` (Enum/String: `user`, `supplier`, `admin`)
  - `address` (String)
  - `document_url` (String - untuk SIUP/NIB, opsional bagi UMKM, wajib bagi Supplier)
  - `status` (String: `pending`, `active`, `suspended` - untuk verifikasi supplier oleh admin)
  - `created_at`, `updated_at`

- **Products:**
  - `id` (UUID)
  - `supplier_id` (Foreign Key ke Users)
  - `name` (String)
  - `category` (String)
  - `price` (Decimal/Float)
  - `stock` (Integer)
  - `description` (Text)
  - `created_at`, `updated_at`

- **Orders:**
  - `id` (UUID)
  - `umkm_id` (Foreign Key ke Users - Pembeli)
  - `supplier_id` (Foreign Key ke Users - Penjual)
  - `total_price` (Decimal/Float)
  - `status` (String: `pending`, `paid`, `processing`, `shipped`, `completed`, `cancelled`)
  - `created_at`, `updated_at`

- **Logs / System Audits (Opsional untuk Admin):**
  - `id`, `action`, `user_id`, `description`, `created_at`

### 3. API Endpoints yang Dibutuhkan

**A. Autentikasi (Public):**
- `POST /api/auth/register` -> Menerima pendaftaran UMKM & Supplier.
- `POST /api/auth/login` -> Menerima email & password, mengembalikan tipe `role` pengguna dan Bearer JWT.

**B. Dashboard UMKM / User (Role: `user`):**
- `GET /api/user/stats` -> Mengembalikan total pesanan bulan ini, jumlah pesanan dikirim, poin, dll.
- `GET /api/user/orders` -> Riwayat pembelian / pesanan.
- `GET /api/products` -> Melihat katalog produk dari semua supplier aktif.
- `POST /api/orders` -> Membuat pesanan baru ke supplier.

**C. Dashboard Supplier (Role: `supplier`):**
- `GET /api/supplier/stats` -> Mengembalikan sisa stok keseluruhan, jumlah pesanan baru, pendapatan, dan rating.
- `GET /api/supplier/products` -> CRUD Manajemen Produk spesifik milik supplier yang sedang login.
- `GET /api/supplier/orders` -> Melihat pesanan masuk dari UMKM.
- `PUT /api/supplier/orders/:id` -> Mengubah status pesanan (misalnya dari `paid` menjadi `shipped`).

**D. Dashboard Admin (Role: `admin`):**
- `GET /api/admin/stats` -> Mengembalikan metrik platform: Total transaksi, pertumbuhan pendapatan, total supplier aktif.
- `GET /api/admin/suppliers` -> Melihat daftar semua supplier (untuk keperluan verifikasi dokumen / NIB).
- `PUT /api/admin/suppliers/:id/verify` -> Menerima atau menolak hak akses supplier masuk ke ekosistem.
- `GET /api/admin/logs` -> Log sistem aktivitas (Log Aktivitas Terbaru).

### 4. Detail Keamanan & Middleware
- Setup **CORS** agar dapat menerima request dari UI Frontend HTML terpisah (izinkan `*` untuk development, gunakan konfigurasi logik untuk *Production*).
- Buatkan **Middleware JWT Mutlak**: `RequireAuth` untuk semua rute `/api/*` kecuali login/register.
- Buatkan **Middleware Role-based**: `RequireRole("admin")`, `RequireRole("supplier")`, dsb., untuk memastikan UMKM tidak bisa mengakses endpoint khusus Admin/Supplier.

Tolong berikan _folder structure_ (contoh standard hexagonal architecture atau MVC standar: `/controllers`, `/models`, `/routes`, `/middlewares`) dan juga kode Golang _bootstrap_ sederhana (main.go dan route handler dasar) sesuai requirement di atas!
