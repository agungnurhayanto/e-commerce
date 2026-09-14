# Application Flow Ecommerce Single Vendor

## 1. Customer Flow

### Melihat Produk

1. Customer membuka website.
2. Customer melihat daftar produk.
3. Customer memilih salah satu produk.
4. Sistem menampilkan detail produk.

### Menambahkan Produk ke Keranjang

1. Customer memilih jumlah produk.
2. Customer menekan tombol "Tambah ke Keranjang".
3. Sistem memeriksa ketersediaan stok.
4. Sistem menyimpan produk ke keranjang.
5. Sistem menampilkan isi keranjang.

### Checkout

1. Customer membuka keranjang.
2. Customer memeriksa produk dan jumlahnya.
3. Customer menekan tombol checkout.
4. Customer mengisi alamat pengiriman.
5. Sistem menghitung total pesanan.
6. Customer mengonfirmasi pesanan.
7. Sistem membuat pesanan baru.

### Melihat Pesanan

1. Customer membuka halaman pesanan.
2. Sistem menampilkan daftar pesanan milik customer.
3. Customer memilih salah satu pesanan.
4. Sistem menampilkan detail pesanan dan statusnya.

---

## 2. Admin Flow

### Login Admin

1. Admin membuka halaman login.
2. Admin memasukkan email dan password.
3. Sistem memeriksa data login.
4. Jika valid, admin masuk ke dashboard.

### Mengelola Kategori

1. Admin membuka halaman kategori.
2. Admin melihat daftar kategori.
3. Admin dapat menambahkan kategori.
4. Admin dapat mengubah kategori.
5. Admin dapat menghapus kategori.

### Mengelola Produk

1. Admin membuka halaman produk.
2. Admin melihat daftar produk.
3. Admin dapat menambahkan produk.
4. Admin dapat mengubah produk.
5. Admin dapat menghapus produk.

### Mengelola Pesanan

1. Admin membuka halaman pesanan.
2. Admin melihat daftar pesanan customer.
3. Admin memilih salah satu pesanan.
4. Admin melihat detail pesanan.
5. Admin mengubah status pesanan.

---

## 3. Status Pesanan

Pesanan memiliki status berikut:

- pending
- processing
- shipped
- completed
- cancelled

Urutan normal:

pending → processing → shipped → completed

Jika pesanan dibatalkan:

pending → cancelled
