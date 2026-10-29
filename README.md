# Ruang Tenang API

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-v1.10-008ECF?style=for-the-badge&logo=go&logoColor=white)](https://gin-gonic.com/)
[![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1?style=for-the-badge&logo=mysql&logoColor=white)](https://www.mysql.com/)
[![Swagger](https://img.shields.io/badge/Swagger-OpenAPI%203.0-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)](http://localhost:8080/swagger/index.html)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)

Backend RESTful API berbasis Golang untuk platform kesehatan mental **Ruang Tenang**. API ini berfungsi sebagai *single source of truth* untuk seluruh data, autentikasi, otorisasi, migrasi database MySQL, kurasi konten, integrasi AI, komunitas, gamifikasi, billing payment gateway, serta analitik B2B organisasi.

---

## Daftar Isi

- [Konteks Ekosistem](#konteks-ekosistem)
- [Fitur dan Domain](#fitur-dan-domain)
- [Prasyarat Sistem](#prasyarat-sistem)
- [Mulai Cepat (Quick Start)](#mulai-cepat-quick-start)
  - [Langkah Instalasi](#langkah-instalasi)
  - [Akun Demo Seeder](#akun-demo-seeder)
- [Daftar Perintah Makefile](#daftar-perintah-makefile)
- [Konfigurasi Environment](#konfigurasi-environment)
- [Arsitektur & Struktur Direktori](#arsitektur--struktur-direktori)
- [Dokumentasi OpenAPI / Swagger](#dokumentasi-openapi--swagger)
- [Troubleshooting & FAQ](#troubleshooting--faq)
- [CI/CD & Deployment](#cicd--deployment)
- [Keamanan & Kontribusi](#keamanan--kontribusi)

---

## Konteks Ekosistem

API ini melayani dua klien utama dalam platform Ruang Tenang:
1. **Frontend Web (`ruang-tenang-web`)**: Melayani pengguna dengan peran *Member*, *Admin/Moderator*, dan *Mitra Organisasi*.
2. **Aplikasi Mobile (`ruang-tenang-mobile`)**: Aplikasi Flutter mobile-first khusus untuk peran *Member*.

Kedua klien berkomunikasi dengan API ini melalui endpoint RESTful di `/api/v1/*`. Role server merupakan batas keamanan (*security boundary*) mutlak—klien tidak dipercaya hanya berdasarkan visibilitas menu di antarmuka.

---

## Fitur dan Domain

- **Autentikasi & Akun**: Registrasi via nomor WhatsApp, OTP Fonnte, login JWT, reset kata sandi, dan proteksi RBAC (*Role-Based Access Control*).
- **Self-Care & Wellness**: Tracker suasana hati (*mood*), jurnal harian (privat/publik), asesmen wellness, dan deteksi krisis/moderasi teks otomatis.
- **AI Mental Health Companion**: Chatbot asisten konseling berbasis DeepSeek dengan manajemen kuota harian, konteks jurnal, dan riwayat percakapan.
- **Komunitas & Konten**: Artikel kesehatan mental berlisensi, kisah nyata (*stories*), forum diskusi interaktif, serta pencarian konten.
- **Relaksasi & Audio**: Katalog musik instrumental terkurasi (lisensi CC BY 4.0) dan manajemen playlist favorit.
- **Gamifikasi & Perjalanan**: Sistem EXP, level peningkatan, perolehan lencana (*badge*), quest harian, leaderboard, serta peta perjalanan self-care.
- **Billing & Transaksi**: Integrasi gateway Duitku Pop (sandbox/production), top-up koin, invoice otomatis, dan pembatasan penggunaan fitur berbayar.
- **Mitra & B2B Organisasi**: Manajemen seat organisasi korporat/kampus, analitik agregat kesehatan mental, SSO, dan audit log.
- **Notifikasi**: Web Push notifications (VAPID) dan siaran pesan (*broadcast*).

---

## Prasyarat Sistem

Pastikan perangkat lokal Anda telah memenuhi spesifikasi berikut:
- **Go**: Versi 1.25 atau lebih baru
- **Database**: MySQL 8.0.13+ atau MariaDB 10.2+ (pastikan service database sudah aktif)
- **Make**: GNU Make untuk menjalankan otomatisasi workflow
- **Tools Tambahan**: `swag` (v1.16.4) dan `migrate` (v4.20.1) — *dapat diinstal otomatis via Makefile*

---

## Mulai Cepat (Quick Start)

### Langkah Instalasi

1. **Salin Environment Template:**
   ```bash
   cp .env.example .env
   ```
   *Edit `.env` dan pastikan konfigurasi `DATABASE_URL` sesuai dengan kredensial MySQL lokal Anda.*

2. **Pasang Tools Generator & Migration:**
   ```bash
   make install-tools
   ```

3. **Jalankan Setup Otomatis (Dependencies, Migrasi, & Seeder):**
   ```bash
   make setup
   ```
   *Perintah ini akan mengunduh Go modules, menjalankan seluruh rantai migrasi SQL, dan mengisi data awal demo.*

4. **Jalankan Server API:**
   ```bash
   # Mode standar
   make run

   # Atau mode hot-reload (memerlukan air)
   make dev
   ```

Server default berjalan pada **`http://localhost:8080`**.
- **Health Check**: `http://localhost:8080/health`
- **Dokumentasi Swagger UI**: `http://localhost:8080/swagger/index.html`

### Akun Demo Seeder

Setelah menjalankan `make seed` atau `make setup`, database akan terisi beberapa akun pengujian lokal:

| Peran (Role) | Email | Password | Kegunaan Pengujian |
| :--- | :--- | :--- | :--- |
| **Admin** | `admin@gmail.com` | `password` | Manajemen user, moderasi forum, kelola krisis, broadcast |
| **Mitra** | `mitra@gmail.com` | `password` | Dashboard B2B organisasi, alokasi seat, laporan agregat |
| **User Premium** | `gading@gmail.com` | `password` | Akses penuh seluruh konten & fitur berbayar |
| **User B2B** | `dery@gmail.com` | `password` | Akses member yang terikat seat organisasi |
| **User Freemium** | `andhika@gmail.com` | `password` | Pengujian kuota fitur gratis & alur upgrade |

> [!CAUTION]
> Akun dan kredensial default di atas **hanya diperuntukkan bagi lingkungan development/demo lokal**. Jangan pernah gunakan akun seeder ini di lingkungan production!

---

## Daftar Perintah Makefile

| Perintah | Deskripsi |
| :--- | :--- |
| `make deps` | Mengunduh dan merapikan Go modules (`go mod tidy`) |
| `make run` | Menjalankan server HTTP API secara langsung |
| `make dev` | Menjalankan server dengan live hot-reload (*Air*) |
| `make build` | Melakukan kompilasi binary server dan seeder ke direktori `bin/` |
| `make test` | Menjalankan unit test (`go test ./...`) |
| `make swagger` | Me-regenerasi file dokumentasi `docs/docs.go` dan `docs/openapi.yaml` |
| `make swagger-check` | Memeriksa apakah terdapat *drift* antara anotasi kode dan file snapshot OpenAPI |
| `make migrate-up` | Menjalankan seluruh migrasi database yang belum dieksekusi |
| `make migrate-down` | Melakukan rollback migrasi terakhir |
| `make migrate-fresh` | Menghapus dan membuat ulang struktur database dari awal *(Destruktif)* |
| `make seed` | Mengisi database dengan data demo/presentasi |
| `make check-port` | Memeriksa apakah port target (default 8080) sedang digunakan |
| `make kill-port` | Menghentikan paksa proses yang sedang memakai port target |
| `make quickstart-check` | Memvalidasi integritas file env, koneksi DB, migrasi, dan health check |
| `make docker-build` | Membangun image Docker untuk API |
| `make docker-run` | Menjalankan container Docker dengan konfigurasi dari file `.env` |

---

## Konfigurasi Environment

Konfigurasi dibaca dari environment variables atau file `.env` melalui `LoadConfig`.

| Kelompok | Variabel | Deskripsi & Contoh Nilai | Wajib/Opsional |
| :--- | :--- | :--- | :--- |
| **Core** | `APP_ENV` | Lingkungan runtime (`development`, `staging`, `production`) | Wajib |
| | `PORT` | Port server HTTP (default: `8080`) | Wajib |
| | `DATABASE_URL` | Format: `mysql://user:pass@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=True` | Wajib |
| | `JWT_SECRET` | Kunci rahasia signing JWT token | Wajib |
| | `CORS_ALLOWED_ORIGINS`| Domain yang diizinkan, dipisah koma (contoh: `http://localhost:3000`) | Wajib |
| **Runtime** | `APP_TIMEZONE` | Zona waktu aplikasi (default: `Asia/Jakarta`) | Opsional |
| | `FRONTEND_URL` | URL basis frontend web (contoh: `http://localhost:3000`) | Opsional |
| | `API_PUBLIC_URL` | URL publik API untuk callback/webhook (contoh: `http://localhost:8080`) | Opsional |
| | `JWT_EXPIRY_HOURS` | Masa berlaku token JWT dalam jam (default: `24`) | Opsional |
| **AI (DeepSeek)**| `DEEPSEEK_API_KEY` | API Key dari DeepSeek untuk fitur Chatbot & Wellness | Wajib jika fitur AI aktif |
| | `DEEPSEEK_BASE_URL`| Endpoint API DeepSeek (default: `https://api.deepseek.com`) | Opsional |
| | `AI_CHAT_MODEL` | Model AI yang digunakan (default: `deepseek-flash`) | Opsional |
| **Billing** | `DUITKU_SANDBOX` | Status mode testing payment (`true` atau `false`) | Wajib |
| | `DUITKU_MERCHANT_CODE`| Kode merchant Duitku Pop | Opsional di lokal |
| | `DUITKU_API_KEY` | Server API key Duitku Pop | Opsional di lokal |
| **WhatsApp** | `FONNTE_TOKEN` | Token API Fonnte untuk verifikasi nomor WhatsApp dan OTP | Wajib jika SMS/OTP aktif |
| **Chat Quota** | `CHAT_DAILY_MESSAGE_LIMIT`| Batas pesan AI per hari (contoh: `100`) | Opsional |
| | `CHAT_QUOTA_RESET_INTERVAL`| Interval reset kuota chat (contoh: `24h`) | Opsional |
| **Web Push** | `VAPID_PUBLIC_KEY` | Public key VAPID untuk Web Push Notifications | Opsional |
| | `VAPID_PRIVATE_KEY`| Private key VAPID | Opsional |
| | `VAPID_CONTACT` | Kontak penanggung jawab (contoh: `mailto:admin@ruang-tenang.site`) | Opsional |

Panduan konfigurasi mendalam Duitku Pop tersedia di [docs/DUITKU_SETUP.md](docs/DUITKU_SETUP.md).  
Atribusi, kredit lisensi musik, dan sumber data seeder tercantum di [docs/SEED_CONTENT.md](docs/SEED_CONTENT.md).

---

## Arsitektur & Struktur Direktori

Repository ini mengadopsi prinsip **Clean Architecture / Domain-Driven Separation**:

```text
ruang-tenang-api/
├── cmd/
│   ├── server/          # Entrypoint utama server HTTP
│   ├── seeder/          # Command pengisian data demo
│   └── migrate/         # Command eksekusi migrasi database
├── internal/
│   ├── features/        # Modul fitur per domain (auth, mood, journal, music, forum, b2b, dll)
│   │   └── <domain>/
│   │       ├── application/      # Use case dan business logic
│   │       ├── infrastructure/   # Akses database repository GORM & adapter eksternal
│   │       └── interface/http/   # Handler HTTP Gin, request DTO, & response mapping
│   ├── router/          # Pendaftaran rute API v1 dan dependency injection
│   ├── middleware/      # Auth JWT, Role check, CORS, Rate Limit, Recovery, Zap Logger
│   ├── model/           # Entitas database GORM lintas domain
│   └── dto/             # Data Transfer Objects bersama
├── migrations/          # File migrasi SQL terurut (didukung golang-migrate)
├── prompts/             # Prompt engineering AI yang di-embed langsung ke binary
├── docs/                # Generated Swagger (docs.go) dan snapshot openapi.yaml
└── pkg/                 # Utilitas helper independen (hash, response formatter, validator)
```

---

## Dokumentasi OpenAPI / Swagger

Anotasi Swagger pada handler HTTP di `internal/features/*` dan `cmd/server/main.go` merupakan *single source of truth* untuk kontrak API.

- **Regenerasi Dokumentasi:**
  ```bash
  make swagger
  ```
- **Validasi Kontrak pada CI:**
  ```bash
  make swagger-check
  ```
  *Perintah ini memastikan tidak ada perubahan kode handler yang belum di-generate ke file `docs/openapi.yaml`.*

---

## Troubleshooting & FAQ

### 1. Port 8080 sudah digunakan (`Address already in use`)
Gunakan utilitas Makefile untuk menghentikan proses sebelumnya atau ubah port:
```bash
make kill-port
# Atau jalankan pada port berbeda:
PORT=8085 make run
```

### 2. Error koneksi database saat migrasi / seeder
Pastikan service MySQL lokal Anda berjalan dan database yang tertera di `DATABASE_URL` sudah dibuat:
```sql
CREATE DATABASE ruang_tenang CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```
Lalu verifikasi koneksi dengan:
```bash
make quickstart-check
```

### 3. Migrasi kotor (*Dirty Migration*)
Jika proses migrasi terhenti di tengah jalan dan berstatus *dirty*, perbaiki migrasi tersebut melalui command migrate atau gunakan `make migrate-fresh` (hanya di development).

---

## CI/CD & Deployment

- **Automated Workflow**: Workflow GitHub Actions `.github/workflows/build-and-deploy.yml` secara otomatis memvalidasi Go 1.25, menjalankan `swag-check`, `go mod tidy`, dan mengompilasi binary.
- **Docker Image**: Mendukung build container mandiri dengan multi-stage build. Flag `RUN_MIGRATE=true` atau `RUN_SEEDER=true` dapat diaktifkan saat kontainer pertama kali di-deploy.

---

## Keamanan & Kontribusi

1. **Jaga Rahasia Kredensial**: Jangan pernah melakukan commit file `.env`, API key DeepSeek, token Fonnte, merchant key Duitku, atau JWT Secret ke repository git.
2. **Pedoman Perubahan Kontrak**: Setiap penambahan endpoint atau modifikasi DTO wajib memperbarui anotasi Swagger dan menyertakan hasil regenerasi `make swagger`.
3. Informasi panduan kerja AI Agent tersedia di `AGENTS.md`.
