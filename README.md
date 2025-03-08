# CLI Web Panel

Sebuah tool CLI untuk mengelola server atau VPS dengan pengalaman sysadmin yang minim. Tool ini menyediakan antarmuka command-line yang mudah digunakan untuk manajemen website, database, dan backup.

## Fitur

- 🌐 **Manajemen Website**
  - Mendukung multiple domain
  - Konfigurasi Caddy otomatis
  - Modul-modul yang bisa dikustomisasi

- 🔒 **Modul Terintegrasi**
  - PHP FastCGI
  - Single Page Application (SPA)
  - Keamanan dasar
  - Manajemen header
  - Pembatasan akses IP

- 📊 **Database Management**
  - Membuat dan menghapus database
  - Manajemen user database
  - Pengaturan hak akses
  - Backup otomatis

- 🔄 **Sistem Backup**
  - Backup harian (incremental)
  - Backup mingguan (full)
  - Rotasi backup otomatis
  - Backup website dan database

- 📊 **Monitoring**
  - CPU usage
  - Memory usage
  - Disk space
  - Status service
  - Log monitoring

## Persyaratan Sistem

- Sistem operasi: Linux (Ubuntu/Debian direkomendasikan)
- RAM minimal: 1GB
- Disk space minimal: 20GB
- Koneksi internet untuk instalasi
- Tidak bisa dijalankan sebagai user root (untuk keamanan)

## Instalasi

1. Download script instalasi:
```bash
curl -O https://raw.githubusercontent.com/doko/cli-webpanel/main/scripts/install.sh
```

2. Berikan permission eksekusi:
```bash
chmod +x install.sh
```

3. Jalankan script instalasi:
```bash
sudo ./install.sh
```

Script akan:
- Menginstall Caddy dan MariaDB
- Menyiapkan struktur direktori
- Mengkonfigurasi permission
- Menginstall CLI tool

## Penggunaan

### Manajemen Website

```bash
# Menambah website baru
webpanel site add domain.com

# Melihat daftar website
webpanel site list

# Menghapus website
webpanel site rm domain.com
```

### Manajemen Modul

```bash
# Melihat modul yang tersedia
webpanel module list-available

# Melihat modul yang terpasang pada domain
webpanel module list domain.com

# Menambah modul ke domain
webpanel module add php domain.com
webpanel module add restrict-access domain.com 192.168.1.0/24

# Menghapus modul dari domain
webpanel module rm php domain.com
```

### Manajemen Database

```bash
# Melihat daftar database
webpanel db list

# Membuat database baru
webpanel db create mydb

# Menghapus database
webpanel db delete mydb

# Manajemen user database
webpanel dbuser list
webpanel dbuser create myuser mypassword
webpanel dbuser delete myuser

# Memberikan akses database ke user
webpanel dbgrant myuser mydb
```

### Manajemen Backup

```bash
# Mengaktifkan backup harian
webpanel backup enable daily domain.com

# Mengaktifkan backup mingguan
webpanel backup enable weekly domain.com

# Menonaktifkan backup
webpanel backup disable daily domain.com
webpanel backup disable weekly domain.com

# Backup database
webpanel dbbackup enable daily mydb
webpanel dbbackup enable weekly mydb
```

### Monitoring

```bash
# Melihat status sistem
webpanel status

# Melihat log
webpanel logs
webpanel logs caddy
webpanel logs mariadb

# Monitoring real-time
webpanel monitor
```

## Struktur Direktori

```
/apps/sites/
    ├── domain1.com/
    │   ├── public/
    │   └── logs/
    └── domain2.com/
        ├── public/
        └── logs/

/usr/local/webpanel/
    ├── config/
    │   ├── global/
    │   └── sites/
    ├── lib/
    │   └── modules/
    └── logs/

/backup/
    ├── daily/
    │   ├── domain.com/
    │   └── dbname/
    └── weekly/
        ├── domain.com/
        └── dbname/
```

## Keamanan

- Tool tidak bisa dijalankan sebagai root user
- SSL/TLS otomatis dengan Caddy
- Firewall configuration
- Enkripsi backup
- Secure configuration storage

## Debugging

Jika terjadi masalah, cek log di:
```bash
tail -f /usr/local/webpanel/logs/webpanel.log
```

## Kontribusi

Kontribusi sangat diterima! Silakan buat pull request untuk:
- Bug fixes
- Fitur baru
- Dokumentasi
- Optimisasi
- Peningkatan keamanan

## Lisensi

MIT License - lihat file [LICENSE](LICENSE) untuk detail lengkap.