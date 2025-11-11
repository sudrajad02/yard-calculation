# Yard Planning Backend Service

## Fitur Utama

*   **Manajemen Rencana Yard:** Mendefinisikan area spesifik dalam block untuk tipe kontainer tertentu (berdasarkan ukuran, tinggi, dan tipe).
*   **Saran Penempatan:** Endpoint untuk mendapatkan saran posisi kontainer berdasarkan rencana yang telah ditentukan.
*   **Penempatan Kontainer:** Endpoint untuk mencatat penempatan kontainer ke posisi tertentu, dengan validasi ketersediaan dan kesesuaian rencana.
*   **Pengambilan Kontainer:** Endpoint untuk mencatat pengambilan kontainer dari lapangan.

## Teknologi yang Digunakan

*   **Bahasa Pemrograman:** [Go (Golang)](https://golang.org/)
*   **Web Framework:** [Fiber](https://gofiber.io/)
*   **ORM (Object-Relational Mapping):** [GORM](https://gorm.io/)
*   **Database:** [PostgreSQL](https://www.postgresql.org/)
*   **Manajemen Environment Variable:** [godotenv](https://github.com/joho/godotenv)

## Struktur Proyek
```
your-project-name/
├── main.go
├── config/
│   └── database.go
├── models/
│   ├── yard.go
│   ├── block.go
│   ├── yard_plan.go
│   └── container.go
├── handlers/
│   └── container_handler.go
├── services/
│   └── container_service.go
├── repositories/
│   └── container_repository.go
├── utils/
│   └── response.go (contoh)
├── .env
└── go.mod
```

## Instalasi dan Persiapan

1.  **Pastikan Go sudah terinstall** (versi minimum 1.18 direkomendasikan).
2.  **Pastikan PostgreSQL sudah terinstall dan berjalan.**
3.  **Clone atau unduh repository ini.**
4.  **Masuk ke direktori proyek:**
    ```bash
    cd path/ke/your-project-name
    ```
5.  **Install dependencies Go:**
    ```bash
    go mod tidy
    ```
6.  **Setup Database:**
    *   Buat database PostgreSQL baru (misalnya `yard_planning_db`).
    *   Sesuaikan nama database, user, dan password di file `.env`.
7.  **Konfigurasi Environment Variables:**
    *   Salin file `.env.sample` (jika ada) atau buat file baru bernama `.env`.
    *   Isi dengan konfigurasi database PostgreSQL kamu:
        ```env
        DB_HOST=localhost
        DB_USER=your_db_user
        DB_PASSWORD=your_db_password
        DB_NAME=your_db_name
        DB_PORT=your_db_port
        ```
        Gantilah `your_db_user`, `your_db_password`, `your_db_name`, dan `your_db_port` dengan nilai yang sesuai dengan setup PostgreSQL kamu.

## Menjalankan Aplikasi

1.  Pastikan konfigurasi database di `.env` sudah benar.
2.  Jalankan perintah berikut dari direktori proyek:
    ```bash
    go run main.go
    ```
3.  Aplikasi akan berjalan secara otomatis melakukan migrasi skema database (membuat tabel `yards`, `blocks`, `yard_plans`, `containers`) jika belum ada, dan mendengarkan permintaan di `http://localhost:3003`.

## API Endpoints

Layanan ini menyediakan RESTful API berikut:

### 1. Get Suggestion Position

Mendapatkan saran posisi untuk meletakkan kontainer berdasarkan rencana.

*   **URL:** `/suggestion`
*   **Method:** `POST`
*   **Content-Type:** `application/json`
*   **Request Body:**
    ```json
    {
      "yard": "YRD1",
      "container_number": "ALFI000001",
      "container_size": 20,
      "container_height": 8.6,
      "container_type": "DRY"
    }
    ```
    *   `yard` (string): ID yard tempat mencari saran.
    *   `container_number` (string): Nomor kontainer.
    *   `container_size` (int): Ukuran kontainer (20 atau 40).
    *   `container_height` (float64): Tinggi kontainer (misalnya 8.6, 9.6).
    *   `container_type` (string): Tipe kontainer (misalnya "DRY", "REEFER").
*   **Response (Success - 200 OK):**
    ```json
    {
      "code": 200,
      "message": "Suggest Container",
      "data": {
        "suggested_position": {
          "block": "LC01",
          "slot": 1,
          "row": 1,
          "tier": 1
        }
      },
      "error": null
    }
    ```
*   **Response (Error - 400/500):**
    ```json
    {
      "code": 500,
      "message": "Error Get Suggest",
      "data": null,
      "error": "detail_error_message"
    }
    ```

### 2. Palace Container

Mencatat bahwa kontainer telah ditempatkan di posisi tertentu.

*   **URL:** `/placement`
*   **Method:** `POST`
*   **Content-Type:** `application/json`
*   **Request Body:**
    ```json
    {
      "yard": "YRD1",
      "container_number": "ALFI000001",
      "block": "LC01",
      "slot": 1,
      "row": 1,
      "tier": 1,
      "container_size": 20,
      "container_height": 8.6,
      "container_type": "DRY"
    }
    ```
    *   Field `yard`, `container_number`, `block`, `slot`, `row`, `tier` harus sesuai dengan posisi yang dituju.
    *   Field `container_size`, `container_height`, `container_type` digunakan untuk validasi kesesuaian rencana.
*   **Response (Success - 200 OK):**
    ```json
    {
      "code": 200,
      "message": "Success",
      "data": null,
      "error": null
    }
    ```
*   **Response (Error - 400/500):**
    ```json
    {
      "code": 500,
      "message": "Error Placing Container",
      "data": null,
      "error": "detail_error_message"
    }
    ```

### 3. Pickup Container

Mencatat bahwa kontainer telah diambil dari lapangan.

*   **URL:** `/pickup`
*   **Method:** `POST`
*   **Content-Type:** `application/json`
*   **Request Body:**
    ```json
    {
      "yard": "YRD1",
      "container_number": "ALFI000001"
    }
    ```
    *   `yard` (string): ID yard tempat kontainer berada.
    *   `container_number` (string): Nomor kontainer yang diambil.
*   **Response (Success - 200 OK):**
    ```json
    {
      "code": 200,
      "message": "Success",
      "data": null,
      "error": null
    }
    ```
*   **Response (Error - 400/500):**
    ```json
    {
      "code": 500,
      "message": "Error Picking Up Container",
      "data": null,
      "error": "detail_error_message"
    }
    ```
