# Task List: Integrasi Supabase & Golang (SupplierHub)

Mengingat proyek ini menggunakan **Golang** sebagai backend dan **Supabase** sebagai *Backend-as-a-Service* (BaaS) / Database, berikut adalah rincian tugas tambahan khusus untuk setup dan integrasi arsitektur tersebut:

## Fase 1: Setup & Konfigurasi Supabase
Supabase menggunakan arsitektur PostgreSQL di belakang layar. Kita akan menghubungkan Golang menggunakan ORM (seperti GORM) atau driver standar Postgres.
- [ ] **Task 1.1:** Buat proyek baru di [Supabase Dashboard](https://supabase.com/).
- [ ] **Task 1.2:** Dapatkan kredensial proyek: `API URL`, `anon/public key`, dan `Database Connection String` (URI PostgreSQL).
- [ ] **Task 1.3:** Setup variabel lingkungan (Environment Variables) di file `.env` Golang Anda:
  ```env
  SUPABASE_URL="https://xyz.supabase.co"
  SUPABASE_KEY="ey..."
  DB_DSN="postgresql://postgres:[PASSWORD]@db.xyz.supabase.co:5432/postgres"
  ```
- [ ] **Task 1.4:** Update file konfigurasi database di backend Go (misalnya `config/database.go`) untuk melakukan koneksi ke DB Supabase menggunakan GORM Postgres Driver.

## Fase 2: Migrasi Model ke Supabase (Database)
- [ ] **Task 2.1:** Sesuaikan file `models/models.go` agar sepenuhnya kompatibel dengan standar Postgres Supabase (seperti tipe data `uuid` bawaan `gen_random_uuid()`).
- [ ] **Task 2.2:** Jalankan `AutoMigrate` pada GORM agar tabel (User, Product, Order, Log) terbuat secara otomatis di dalam database Supabase.
- [ ] **Task 2.3:** (Opsional) Mengaktifkan fitur **Row Level Security (RLS)** melalui SQL Editor di Dashboard Supabase untuk mengamankan data transaksi tiap supplier/UMKM secara absolut dari sisi database.

## Fase 3: Integrasi Supabase Auth (Autentikasi Pengguna)
Karena Anda menggunakan Supabase, Anda bisa menggunakan sistem autentikasi bawaan Supabase yang jauh lebih aman daripada membuat JWT sendiri.
- [ ] **Task 3.1:** Install *Supabase Go Client* resmi: `go get github.com/supabase-community/supabase-go`
- [ ] **Task 3.2:** Modifikasi `controllers/auth_controller.go` (Fungsi Login & Register) agar menembak API pendaftaran Supabase (`supabase.Auth.SignUp` dan `SignIn`).
- [ ] **Task 3.3:** Modifikasi `middlewares/auth.go` untuk memvalidasi token JWT bawaan Supabase alih-alih token lokal, dan mengekstrak `user_id` dari JWT tersebut.

## Fase 4: Integrasi Supabase Storage (Untuk File Upload)
Pada *flow* pendaftaran, Supplier wajib mengunggah dokumen legalitas (SIUP/Akta). File ini sangat ideal disimpan di **Supabase Storage**.
- [ ] **Task 4.1:** Buat *bucket* baru di Supabase Storage bernama `supplier_documents` dan set privasinya.
- [ ] **Task 4.2:** Buat service baru `services/storage.go` di Golang.
- [ ] **Task 4.3:** Buat fungsi Go untuk mengunggah file via REST API Supabase (menggunakan endpoint Storage) saat terjadi proses submit pendaftaran supplier, lalu simpan URL yang dihasilkan ke field `document_url` di tabel User/Supplier.

## Fase 5: Implementasi Real-time Subscriptions (Opsional / Nilai Plus)
Supabase memiliki fitur Real-time (Websocket) bawaan. Ini sangat berguna untuk fitur "Notifikasi Pembayaran / Konfirmasi Stok".
- [ ] **Task 5.1:** Set-up fitur *Realtime* pada tabel `Order` di dashboard Supabase.
- [ ] **Task 5.2:** Buat endpoint WebSocket di Golang atau langsung dengarkan perubahan dari klien (Frontend JavaScript) menggunakan *Supabase JS Client* saat status order UMKM diubah oleh Supplier.
