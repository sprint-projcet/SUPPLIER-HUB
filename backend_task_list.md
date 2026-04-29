# Daftar Pekerjaan (Task List) Pengembangan Backend SupplierHub

Berdasarkan analisis arsitektur (tabel spesifikasi API & flowchart), sistem SupplierHub mengelola hubungan B2B antara UMKM dan pemasok bahan baku, di mana SupplierHub tidak memproses pembayaran secara langsung melainkan mengandalkan **SmartBank** (untuk pemrosesan finansial) dan **LogistiKita** (opsional untuk pengiriman).

Berikut adalah daftar pekerjaan untuk mengimplementasikan fungsionalitas tersebut pada *backend* (Go + Gin):

## Fase 1: Setup & Integrasi Sistem Eksternal
*SupplierHub membutuhkan modul HTTP Client untuk berkomunikasi dengan sistem inti lainnya.*
- [ ] **Task 1.1:** Buat direktori `services/` di dalam root backend.
- [ ] **Task 1.2:** Buat file `services/smartbank.go`. Implementasikan fungsi HTTP Client untuk menembak API `POST /supplier/pay` ke *API Gateway / SmartBank* (menangani pembuatan payment request).
- [ ] **Task 1.3:** Buat file `services/logistikita.go`. Implementasikan fungsi integrasi `POST /logistics/pay` untuk menghitung dan mencatat ongkos kirim.
- [ ] **Task 1.4:** Modifikasi file `models/models.go` untuk menyiapkan struktur database yang bisa menyimpan data log transaksi ke pihak ketiga (misalnya entitas `PaymentLog` atau penambahan atribut respon SmartBank pada struct `Order`).

## Fase 2: Implementasi Endpoint Spesifik (Berdasarkan Tabel)
*Menyesuaikan endpoint routing di file `routes/routes.go` agar presisi dengan desain yang diminta.*
- [ ] **Task 2.1: Fitur Manajemen Bahan Baku**
  - **Endpoint:** `/supplierhub/manajemen_bahan_baku`
  - **Deskripsi:** Memungkinkan supplier menambah/mengubah stok bahan baku.
  - **File Target:** `controllers/supplier_controller.go`
- [ ] **Task 2.2: Fitur Biaya Layanan Supplier**
  - **Endpoint:** `/supplierhub/biaya_layanan_supplier`
  - **Deskripsi:** Modul untuk menghitung potongan *fee* platform per transaksi supplier sebelum dikirim ke SmartBank.
  - **File Target:** `controllers/supplier_controller.go`

## Fase 3: Alur Inti (Order & Pembayaran B2B)
*Menerjemahkan alur dari UMKM memesan bahan baku hingga proses request pembayaran berjalan.*
- [ ] **Task 3.1: Fitur Order Bahan**
  - **Endpoint:** `/supplierhub/order_bahan`
  - **Deskripsi:** Menerima *request* pesanan dari UMKM.
  - **Flow:** Validasi input *item & qty* -> Cek ketersediaan stok & harga -> Hitung subtotal -> Hitung biaya layanan tambahan (Task 2.2) -> Generate `Payment Request` via fungsi di `services/smartbank.go` -> Simpan status `Order` ke DB.
  - **File Target:** `controllers/user_controller.go`
- [ ] **Task 3.2: Fitur Pembayaran**
  - **Endpoint:** `/supplierhub/pembayaran`
  - **Deskripsi:** Menangani respons atau mekanisme integrasi dari sistem SmartBank setelah pemrosesan debit/kredit dan pajak selesai dilakukan di *Core SmartBank*.
  - **Flow:** Validasi respon SmartBank -> Ubah status `Order` (misal dari Pending ke Paid) -> Catat ke dalam Log.
  - **File Target:** `controllers/user_controller.go`

## Fase 4: Fulfillment oleh Supplier
- [ ] **Task 4.1: Fitur Konfirmasi Stok**
  - **Endpoint:** `/supplierhub/konfirmasi_stok`
  - **Deskripsi:** Supplier memverifikasi ketersediaan dan menyiapkan fisik bahan baku yang sudah dibayar oleh UMKM.
  - **Flow:** Ubah status pesanan -> Memotong kuantitas stok produk di database -> (Jika terintegrasi otomatis) Memanggil request pengiriman via `services/logistikita.go`.
  - **File Target:** `controllers/supplier_controller.go`

## Fase 5: Keamanan & Monitoring (API Gateway Standard)
- [ ] **Task 5.1:** Update fungsi validasi token JWT di `middlewares/auth.go` agar strukturnya kompatibel dengan standar validasi di level `API Gateway / Integrator`.
- [ ] **Task 5.2:** Implementasi pembuatan *log* histori aktivitas sistematis pada database setiap kali terjadi mutasi stok atau perubahan status pesanan. Log ini nantinya dapat diekspor atau dibaca oleh pihak eksternal untuk analisis *UMKM Insight*.
