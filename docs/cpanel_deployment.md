# Panduan Deployment Golang API (`ruang-tenang-api`) ke cPanel Shared Hosting

Dokumen ini berisi panduan teknis langkah demi langkah untuk melakukan deployment backend API berbasis **Golang** (`ruang-tenang-api`) ke shared hosting cPanel menggunakan arsitektur **Phusion Passenger Wrapper (Node.js)** dan **cPanel Terminal (SSH)**.

> [!IMPORTANT]
> Seluruh proses kompilasi biner Golang dilakukan **secara lokal** (*cross-compilation* untuk Linux amd64) menggunakan script `scripts/deploy_cpanel.sh`. Shared hosting cPanel tidak memerlukan instalasi Go compiler maupun tool migrasi eksternal.

---

## 1. Arsitektur Deployment di Shared Hosting

Shared hosting cPanel umumnya tidak menyediakan layanan *daemon supervisor* kustom (seperti `systemd` atau `supervisor`) untuk pengguna reguler. Oleh karena itu, backend dijalankan menggunakan arsitektur **Reverse Proxy Wrapper**:

```mermaid
flowchart LR
    Client["Pengunjung / Frontend"] -->|HTTPS :443| LiteSpeed["Apache / LiteSpeed Web Server"]
    LiteSpeed -->|Passenger Socket| NodeWrapper["Node.js Reverse Proxy (app.js)"]
    NodeWrapper -->|Internal HTTP :3001| GoBinary["Biner Statis Golang (app-main)"]
    GoBinary -->|TCP :3306| MySQL[("MySQL cPanel (127.0.0.1)")]
```

1. **Phusion Passenger** di cPanel mengelola siklus hidup server Node.js melalui entrypoint `app.js`.
2. **`app.js`** men-*spawn* biner Golang `app-main` pada port internal (`3001`), kemudian meneruskan seluruh request HTTP masuk ke biner Golang.
3. Biner Golang terisolasi secara aman dan mengakses database MySQL lokal `127.0.0.1:3306`.

---

## 2. Prasyarat Lingkungan

### A. Lingkungan Lokal (Komputer Pengembang)
- **Go**: Versi 1.22+ atau yang kompatibel dengan `go.mod`.
- **Tooling**: `bash`, `zip`, Git, Make (opsional).

### B. Lingkungan Server (cPanel Shared Hosting)
- Fitur **Setup Node.js App** (CloudLinux NodeJS Selector / Phusion Passenger).
- Akses **Terminal cPanel** atau koneksi SSH.
- Akses **File Manager** cPanel.
- Subdomain publik: `api.ruang-tenang.my.id` (A Record mengarah ke server hosting dan SSL aktif).
- **Database MySQL / MariaDB**: Dibuat melalui menu **MySQL Databases** di cPanel (beserta user dan password yang memiliki akses `ALL PRIVILEGES`).

---

## 3. Proses Cross-Compile & Perakitan Bundle di Lokal

Jalankan perintah build dari root direktori `ruang-tenang-api`:

```bash
make deploy-cpanel
# atau jalankan script langsung:
bash scripts/deploy_cpanel.sh
```

### Rangkaian Otomasi Script:
1. **Preflight**: Memverifikasi ketersediaan Go compiler, zip, template `deployment/app.js`, folder `migrations/`, dan `storage/`.
2. **Cross-Compilation (Linux amd64)**:
   Membangun tiga biner statis dengan `CGO_ENABLED=0 GOOS=linux GOARCH=amd64`:
   - `app-main` — Server HTTP REST API utama.
   - `app-migrate` — CLI pengelola skema migrasi database.
   - `app-seeder` — CLI pengisi data awal (*presentation seeding*).
3. **Perakitan Bundle**:
   - Menyalin ketiga biner executable.
   - Menyalin template reverse proxy `deployment/app.js` menjadi `app.js`.
   - Menyalin folder `migrations/` dan `storage/`.
   - Menyalin folder `uploads/` (mempertahankan file yang sudah diupload user maupun seed).
   - Memberikan hak akses eksekusi (`chmod +x`) pada seluruh biner.
4. **Pengemasan Zip**: Menghasilkan berkas arsip `upload-api.zip` yang membungkus folder `ruang-tenang-api/`.

### Struktur Berkas di Dalam `upload-api.zip`:
```text
upload-api.zip
└── ruang-tenang-api/      <-- Application Root saat diekstrak di server
    ├── app-main           <-- Biner API Golang utama (statis Linux amd64)
    ├── app-migrate        <-- Biner CLI migrasi database
    ├── app-seeder         <-- Biner CLI seeder data
    ├── app.js             <-- Node.js reverse proxy wrapper
    ├── migrations/        <-- Berkas skrip SQL migrasi (.up.sql & .down.sql)
    ├── storage/           <-- Aset bawaan (badge level, ikon kategori)
    └── uploads/           <-- Berkas user upload & media publik
```

---

## 4. Konfigurasi Awal Node.js App di cPanel

*(Lakukan langkah ini satu kali pada setup pertama)*

1. Masuk ke **cPanel** -> Buka menu **Setup Node.js App**.
2. Klik tombol **Create Application**:
   - **Node.js Version**: Pilih versi LTS (misal `20.x` atau `22.x`).
   - **Application Mode**: `Production`.
   - **Application Root**: `ruang-tenang-api` (relatif terhadap direktori home user `~`).
   - **Application URL**: Pilih subdomain `api.ruang-tenang.my.id`.
   - **Application Startup File**: `app.js`.
3. Klik tombol **Create** / **Save**.

---

## 5. Pengunggahan, Pemasangan & Konfigurasi Server

### A. Unggah Berkas Bundle
Unggah file `upload-api.zip` ke direktori home (`/home/username/`) menggunakan File Manager atau cPanel UAPI / FTP.

### B. Ekstraksi & Setup Permissions via SSH Terminal
Buka **Terminal cPanel** (SSH) dan jalankan:

```bash
cd ~
# Ekstrak file zip (overwrite berkas lama)
unzip -o -q upload-api.zip -d ~

# Masuk ke direktori aplikasi
cd ~/ruang-tenang-api

# Berikan izin eksekusi pada biner dan izin tulis pada folder asset
chmod +x app-main app-migrate app-seeder
chmod -R 775 uploads storage

# Bersihkan file zip sementara untuk menghemat kuota hosting
rm -f ~/upload-api.zip
```

> [!TIP]
> **Gunakan Ekstraksi Terminal SSH**
> Mengekstrak via File Manager cPanel terkadang diblokir oleh sistem dengan error `Permission denied` jika folder target sudah ada. Perintah `unzip -o -q upload-api.zip -d ~` melalui SSH Terminal terbukti 100% andal menimpa berkas lama.

### C. Konfigurasi File `.env` Produksi Server
Buat atau edit berkas `~/ruang-tenang-api/.env`:

```env
# Mode Aplikasi & Zona Waktu (Wajib)
APP_ENV=production
APP_TIMEZONE=Asia/Jakarta

# Port Internal yang Dilayani Golang (di-proxy oleh app.js)
INTERNAL_PORT=3001

# Konfigurasi Database MySQL cPanel (Gunakan 127.0.0.1)
DATABASE_URL=mysql://db_user:db_password@tcp(127.0.0.1:3306)/db_name?charset=utf8mb4&parseTime=True

# Kunci Rahasia JWT (Wajib acak & kuat)
JWT_SECRET=buat-string-acak-panjang-dan-unik-di-sini
JWT_EXPIRY_HOURS=24

# CORS & Domain Origin (Wajib diisi domain frontend)
CORS_ALLOWED_ORIGINS=https://ruang-tenang.my.id,https://www.ruang-tenang.my.id
FRONTEND_URL=https://ruang-tenang.my.id
API_PUBLIC_URL=https://api.ruang-tenang.my.id

# Integrasi AI (DeepSeek - Opsional untuk fitur Chat/Jurnal)
DEEPSEEK_API_KEY=your_deepseek_api_key
DEEPSEEK_BASE_URL=https://api.deepseek.com

# Integrasi Payment Gateway (Duitku - Opsional)
DUITKU_SANDBOX=false
DUITKU_MERCHANT_CODE=
DUITKU_API_KEY=

# Integrasi WhatsApp Gateway (Fonnte - Wajib diisi agar verifikasi OTP dan reset password berfungsi)
# Dapatkan token dari https://fonnte.com (pastikan status device: Connect)
FONNTE_TOKEN=your_fonnte_device_token
```

> [!IMPORTANT]
> **Penyebab Utama WhatsApp Gagal Terkirim di Produksi:**
> 1. **`FONNTE_TOKEN` Kosong di `.env`**: Jika variabel `FONNTE_TOKEN` di server produksi tidak diisi, backend tidak dapat menginisialisasi client WhatsApp dan pengiriman OTP akan ditolak dengan error *"pengiriman WhatsApp belum dikonfigurasi"*.
> 2. **Perangkat WhatsApp Terputus (Disconnected)**: Pastikan perangkat Anda di dashboard Fonnte berada dalam status **Connect**. Uji status perangkat langsung dari Terminal server menggunakan perintah:
>    ```bash
>    curl -s -X POST https://api.fonnte.com/device -H "Authorization: TOKEN_FONNTE_ANDA"
>    ```
>    Jika respons memuat `"device_status":"disconnect"`, scan ulang QR code di dashboard Fonnte.
> 3. **Format Nomor WhatsApp**: Backend secara otomatis menormalisasi nomor ke format `628xxxxxxxxxx` (tanpa tanda `+`) agar sesuai dengan format yang diterima Fonnte API.

### D. Pengamanan Berkas Biner (`.htaccess`)
Tambahkan aturan berikut ke dalam `~/ruang-tenang-api/.htaccess` agar biner aplikasi tidak dapat diunduh publik:

```apache
# Blokir akses langsung ke berkas biner dan source code
<FilesMatch "^(app-main|app-migrate|app-seeder|app\.js|\.env)$">
  Require all denied
</FilesMatch>
Options -Indexes
```

> [!CAUTION]
> **JANGAN HAPUS** blok konfigurasi Passenger bawaan cPanel (`# DO NOT REMOVE. CLOUDLINUX PASSENGER CONFIGURATION`) di dalam file `.htaccess`. Cukup letakkan aturan blokir di atas blok Passenger.

---

## 6. Panduan CLI Migrasi & Seeder Database

Biner CLI `app-migrate` dan `app-seeder` membaca konfigurasi database langsung dari berkas `.env` di direktori yang sama (`DATABASE_URL`).

### A. Migrasi Database (`./app-migrate`)

Jalankan perintah ini dari dalam folder `~/ruang-tenang-api`:

1. **Menerapkan Migrasi Baru (`migrate up`)**:
   ```bash
   ./app-migrate up
   # atau cukup jalankan:
   ./app-migrate
   ```
   *Fungsi*: Menerapkan seluruh berkas skrip SQL baru dari folder `migrations/`. Perintah ini aman dijalankan ulang (*idempotent*).

2. **Membatalkan Migrasi (`migrate down`)**:
   ```bash
   # Rollback 1 step migrasi terakhir:
   ./app-migrate down
   # Rollback sejumlah N step:
   ./app-migrate down 2
   # Rollback SELURUH tabel:
   ./app-migrate down all
   ```

3. **Reset Total Database & Migrasi Ulang dari Nol (`migrate fresh`)**:
   ```bash
   ./app-migrate fresh
   ```
   *Fungsi*: Me-reset database secara total dengan melakukan `DROP` pada seluruh tabel yang ada, kemudian menerapkan seluruh berkas migrasi dari versi pertama hingga terakhir.
   > [!WARNING]
   > `migrate fresh` akan menghapus seluruh data yang ada secara permanen. Hanya gunakan perintah ini untuk setup awal atau di lingkungan testing/demo.

4. **Memperbaiki Status Dirty State (`migrate force`)**:
   Jika migrasi terhenti di tengah jalan karena query gagal atau koneksi putus, status tabel menjadi `dirty: true`. Perbaiki dengan memaksa versi skema kembali ke versi sukses terakhir (misal versi 16):
   ```bash
   ./app-migrate force 16
   ./app-migrate up
   ```

5. **Mengecek Versi Migrasi Saat Ini**:
   ```bash
   ./app-migrate version
   ```

### B. Seeder Data Awal (`./app-seeder`)

Setelah tabel database dibuat (via `up` atau `fresh`), isi database dengan katalog dan data presentasi:

```bash
./app-seeder
```

> [!IMPORTANT]
> **Mekanisme Penyalinan Gambar oleh Seeder:**
> Selain mengisi data baris MySQL, biner `./app-seeder` secara otomatis menyalin berkas fisik dari `storage/images/` ke `uploads/images/` (artikel, kategori, musik, cover cerita, avatar demo). Oleh karena itu, pastikan folder `uploads/` memiliki izin tulis (`chmod -R 775 uploads`).

---

## 7. Me-restart & Memverifikasi Layanan API

1. Pemicu restart Phusion Passenger:
   ```bash
   touch ~/ruang-tenang-api/tmp/restart.txt
   ```
2. **Uji Coba Endpoint Health Check**:
   ```bash
   curl -s https://api.ruang-tenang.my.id/health
   ```
   Respon yang diharapkan:
   ```json
   {"service":"ruang-tenang-api","status":"ok","version":"1.0.0"}
   ```
3. **Uji Coba Akses Gambar Publik**:
   ```bash
   curl -s -I https://api.ruang-tenang.my.id/uploads/images/story-community-support.webp
   ```
   *Respon harus `HTTP/2 200` dengan ukuran konten non-nol.*

---

## 8. Panduan Penyelesaian Masalah (Troubleshooting)

### A. Respon 502 Bad Gateway Saat Server Baru Direstart
- **Penyebab**:
  Wrapper `app.js` memiliki jeda startup (600ms) untuk membersihkan proses lama sebelum men-*spawn* `app-main`. Jika request datang tepat pada detik pertama restart, `app.js` mengembalikan JSON 502 sementara.
- **Penyebab Lanjutan**:
  Jika 502 menetap, kemungkinan biner `app-main` mengalami crash saat boot (misalnya koneksi database MySQL salah atau `.env` tidak lengkap).
- **Solusi**:
  1. Jalankan pengujian manual di Terminal cPanel:
     ```bash
     cd ~/ruang-tenang-api && ./app-main
     ```
  2. Perhatikan log error di layar (misal: `Failed to connect to database`). Perbaiki kredensial di file `.env`.
  3. Tekan `Ctrl + C` setelah pengujian manual selesai.

### B. Error `bind: address already in use` pada Port 3001
- **Penyebab**:
  Proses `app-main` lama masih berjalan di latar belakang dan belum dilepas oleh sistem operasi.
- **Solusi**:
  Bunuh proses yatim di Terminal cPanel lalu restart aplikasi:
  ```bash
  pkill -9 -f app-main
  touch ~/ruang-tenang-api/tmp/restart.txt
  ```

### C. Ekstraksi Zip Gagal (`Permission Denied` di File Manager cPanel)
- **Penyebab**:
  Proses daemon File Manager cPanel terkadang tidak memiliki izin menimpa folder yang sebelumnya dibuat oleh user shell.
- **Solusi**:
  Gunakan Terminal SSH:
  ```bash
  unzip -o -q ~/upload-api.zip -d ~
  ```

### D. Gambar Upload atau Thumbnail Mengembalikan 404
- **Penyebab**:
  File fisik belum disalin ke `uploads/images/` atau permission folder `uploads/` tidak mengizinkan akses baca/tulis.
- **Solusi**:
  Jalankan perintah berikut di direktori `~/ruang-tenang-api`:
  ```bash
  chmod -R 775 uploads storage
  ./app-seeder
  ```

### E. Domain Menampilkan "Index of /" (Directory Listing)
- **Penyebab**:
  Konfigurasi Passenger pada berkas `.htaccess` hilang atau dinonaktifkan oleh web server.
- **Solusi**:
  Buka cPanel -> **Setup Node.js App** -> Klik tombol **RESTART** pada aplikasi terkait, atau hapus dan buat ulang entri aplikasi agar cPanel menulis ulang konfigurasi Passenger pada `.htaccess`.
