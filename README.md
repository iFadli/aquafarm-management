# aquafarm-management
Proyek ini dibuat untuk menjawab Soal Coding Test.

Proyek ini dibuat menggunakan salah satu contoh dari [repository `awesome-compose`](https://github.com/docker/awesome-compose/) dengan stack `nginx-golang-mariadb` yang bisa anda lihat [di sini](https://github.com/docker/awesome-compose/tree/master/nginx-golang-mysql).

Teknologi yang digunakan pada Proyek ini adalah :
- Nginx : Web Server (proxy)
- MariaDB : SQL Database Server (MySQL)
- GoLang : Bahasa Pemrograman Backend (IDE)
- Gin : Framework yang menggunakan GoLang
- Docker : Container
- Git : Sistem Pengontrol Versi (Kode)
- Swagger : Dokumentasi API

## Bagaimana Cara Menjalankan Proyek Ini?
Ada beberapa langkah mudah untuk menjalankan Proyek ini namun pastikan bahwa Komputer anda telah terinstall Docker dan Git.

### 1. Clone Proyek
```
$ mkdir PROJECT && cd PROJECT
$ git clone https://github.com/iFadli/aquafarm-management.git
... # tunggu hingga selesai # ...
Cloning into 'aquafarm-management'...
remote: Enumerating objects: xxx, done.
remote: Counting objects: 100% (xxx/xxx), done.
remote: Compressing objects: 100% (xxx/xxx), done.
remote: Total xxx (delta xx), reused xxx (delta xx), pack-reused x
Receiving objects: 100% (xxx/xxx), xx.xx KiB | xxx.xx KiB/s, done.
Resolving deltas: 100% (xx/xx), done.
$ cd aquafarm-management
```

### 2. * Atur File .env Jika Ingin Merubah DB Config
```
$ nano .env

DB_HOST=db-hub.docker
DB_NAME=aquafarm
DB_USER=root
DB_PASSWORD=BaGuViX91oo
```

### 3. Jalankan Proyek ini dengan docker-compose
```
$ docker-compose up -d
... # tunggu hingga selesai # ...
[+] Running 5/5
 ⠿ Network aquafarm-management-docker_default      Created                                                                                                                      0.0s
 ⠿ Volume "aquafarm-management-docker_db-data"     Created                                                                                                                      0.0s
 ⠿ Container aquafarm-management-docker-db-1       Started                                                                                                                      0.4s
 ⠿ Container aquafarm-management-docker-backend-1  Started                                                                                                                      0.5s
 ⠿ Container aquafarm-management-docker-proxy-1    Started                                                                                                                      0.7s
```

### 4. Akses Proyek sesuai Kebutuhan
Pada pengaturan Docker Proyek ini, secara default akan meng-expose Port 2 Service yang digunakan; Yakni, Database (MariaDB : 3306) dan Web Server (Nginx : 8080).

Jika ingin menghubungkan Database dengan Tool Database Manager seperti DBeaver, anda dapat menyesuaikan konfirugasi dengan File `.env`.
#### !! DBeaver
Berikut langkah-langkah konfigurasi DBeaver :

>1. Buka `DBeaver`.
>2. Klik ikon `Connect to a Database` di Pojok-Kiri-Atas.
>3. Pada Kategori `Popular`, pilih `MariaDB` lalu klik Next.
>4. Isikan Konfigurasi sesuai dengan File `.env` pada Proyek ini.
>5. Jika berhasil, akan ada Ikon centang hijau pada daftar Database di sebelah kiri.
#### !! Swagger
Selain pengaturan Database Manager dari luar Docker yang dapat mengakses ke service Database, di Proyek ini juga disematkan teknologi Swagger untuk memudahkan calon pengguna dalam pengoprasian API ini dengan membaca Dokumentasi API yang cukup Interaktif.

Berikut cara mengaksesya :
>1. Pastikan Instance pada Docker yang tadi kita `compose-docker up -d` statusnya Healthy.
>2. Buka Browser dan akses [`http://localhost:8080/swagger/index.html`](http://localhost:8080/swagger/index.html).
>3. Ketika Swagger berhasil dibuka, pengguna dapat mencoba API yang tertera pada halaman. Perlu di ingat, Proyek ini membutuhkan Authorization ApiKey pada Header disetiap Requestnya. Secara Default ketika proyek ini pertama kali dijalankan, sistem akan mengeksekusi penambahakan ApiKey yang dapat digunakan untuk Percobaan. KeyValue tersebut adalah `ini_dia_si_jali_jali`.
>4. Pada halaman Swagger, anda dapat mengatur Value dari ApiKey yang akan digunakan pada Header dengan cara Klik tombol `Authorize` di sebelah kanan yang berwarna border Hijau. Lalu isikan Valuenya dan Klik `Authorize` hingga tombol berubah menjadi `Logout`. Lalu klik tombol `Close` untuk menutup tampilan modal.
>5. Pada tahap ini, anda telah dapat mengakses API melalui Swagger.

## Daftar API

Sesuai dengan kebutuhan dari Soal Proyek ini, berikut daftar API-nya :
```
[GIN-debug] GET    /v1/farm/
[GIN-debug] GET    /v1/farm/:farm_id
[GIN-debug] POST   /v1/farm/
[GIN-debug] PUT    /v1/farm/
[GIN-debug] DELETE /v1/farm/:farm_id

[GIN-debug] GET    /v1/pond/
[GIN-debug] GET    /v1/pond/:farm_id/:pond_id
[GIN-debug] POST   /v1/pond/
[GIN-debug] PUT    /v1/pond/
[GIN-debug] DELETE /v1/pond/:farm_id/:pond_id

[GIN-debug] GET    /v1/logs/
[GIN-debug] GET    /v1/logs/statistics
```
```
[GIN-debug] GET    /swagger/*                   Path untuk membuka Swagger.
```

---
### Daftar Port Aktif
```
8080   Service Nginx
3306   Service MariaDB
```

Lakukan perintah berikut di Terminal untuk Menonaktifkan Proyek :
```
$ docker-compose down --volumes
```

