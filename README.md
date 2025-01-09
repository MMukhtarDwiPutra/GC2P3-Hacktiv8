[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/7NDPbU0p)
# Graded Challenge 2 - P3

Graded Challenge ini dibuat guna mengevaluasi pembelajaran pada Hacktiv8 Program Fulltime Golang khususnya pada pembelajaran grpc dan implementasi auth dengan Tema :

`REST API menggunakan gRPC dengan Autentikasi untuk Sistem Manajemen Buku`

## **Assignment Objective**  
Graded Challenge 2 ini dibuat untuk mengevaluasi pemahaman Anda terkait konsep gRPC sebagai berikut:
- Mampu memahami konsep gRPC.
- Mampu membuat service dengan gRPC.
- Mampu mengimplementasikan Unit Testing
- Mampu melakukan deployment ke Google Cloud menggunakan Docker image.

## **Assignment Directions**

Anda akan membangun sebuah Sistem Manajemen Buku dengan REST API berbasis gRPC dan autentikasi. Sistem ini mengelola data buku, pengguna, dan memungkinkan pengguna untuk menambahkan, menghapus, atau memperbarui buku yang dimiliki.

Tugas Anda meliputi:
1. **Buat file proto** untuk mendefinisikan layanan gRPC dengan pesan (message) yang sesuai.
2. **Implementasikan server gRPC** untuk layanan Sistem Manajemen Buku dalam Golang.
3. **Tambahkan autentikasi menggunakan token JWT** untuk memastikan hanya pengguna terdaftar yang dapat mengelola buku.
4. **Dokumentasikan seluruh service menggunakan Swagger.**
5. **Lakukan deployment ke Google Cloud** menggunakan Docker image.
6. **Lakukan pengujian menggunakan unit testing** untuk minimal 3 fungsi.
7. **Implementasikan fitur job scheduling** untuk secara otomatis memperbarui status buku yang sudah kadaluarsa atau belum dikembalikan sesuai batas waktu.
8. **Lakukan pengujian/testing pada minimal 3 function dari keseluruhan project

## **Database Schema**
Database bebas (SQL atau NoSQL).

### Tabel Users:
- ID (UUID): Identifikasi unik pengguna.
- Username (String): Nama pengguna.
- Password (String): Kata sandi yang di-hash.

### Tabel Books:
- ID (UUID): Identifikasi unik buku.
- Title (String): Judul buku.
- Author (String): Penulis buku.
- PublishedDate (Timestamp): Tanggal terbit buku.
- Status (String): Status buku (Available, Borrowed).
- UserID (UUID): ID pengguna yang meminjam buku.

### Tabel BorrowedBooks:
- ID (UUID): Identifikasi unik peminjaman.
- BookID (UUID): ID buku yang dipinjam.
- UserID (UUID): ID pengguna yang meminjam.
- BorrowedDate (Timestamp): Tanggal peminjaman.
- ReturnDate (Timestamp): Tanggal pengembalian (null jika belum dikembalikan).

## **Fitur Wajib:**
1. **Autentikasi menggunakan JWT**: Hanya pengguna yang terdaftar yang dapat mengakses dan memanipulasi data buku.
2. **Job Scheduling**: Implementasikan fitur job scheduling yang memperbarui status buku yang belum dikembalikan setelah melewati tenggat waktu peminjaman.
3. **Deployment ke GCP**: Setelah aplikasi berjalan dengan baik secara lokal, Anda harus melakukan deployment ke Google Cloud dan membuat aplikasi dapat diakses secara publik.
4. **Unit Test**: Buat unit test minimal untuk 3 fungsi utama dalam aplikasi Anda.
5. **Docker**: Kontainerisasi aplikasi Anda menggunakan Docker. Pastikan menyediakan Dockerfile dan dokumentasi singkat mengenai bagaimana cara menjalankan aplikasi ini menggunakan Docker.

## **Expected Output:**
1. **File Proto**: Berisi definisi layanan gRPC yang Anda bangun.
2. **Server gRPC**: Berfungsi sebagai backend service untuk aplikasi Sistem Manajemen Buku.
3. **Endpoint gRPC**: Harus diakses melalui URL Google Cloud, contoh: http://url-google-cloud.com/
4. **Swagger Documentation**: Tampilkan dokumentasi service gRPC dan API Anda menggunakan Swagger.
5. **Job Scheduling**: Tugas otomatis untuk memperbarui status buku secara berkala jika belum dikembalikan.
6. **Unit Test**: Pengujian untuk memastikan beberapa fungsi penting berjalan dengan baik.
7. **Deployment Documentation**: Sertakan langkah-langkah deployment ke Google Cloud menggunakan Docker image.


### Additional Notes
Total Points : 100

Deadline : Diinformasikan oleh instruktur saat briefing GC. Keterlambatan pengumpulan tugas mengakibatkan skor GC 2 menjadi 0.

Informasi yang tidak dicantumkan pada file ini harap dipastikan/ditanyakan kembali kepada instruktur. Kesalahan asumsi dari peserta mungkin akan menyebabkan kesalahan pemahaman requirement dan mengakibatkan pengurangan nilai.

### Deployment Notes
- Deployed url: _________ (isi dengan url hasil deployment anda)
