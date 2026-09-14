# ERD Ecommerce Single Vendor

## 1. Daftar Entitas

Aplikasi menggunakan entitas berikut:

- users
- categories
- products
- product_images
- carts
- cart_items
- orders
- order_items
- payments
- addresses

---

## 2. Relasi Antar Entitas

- Satu user dapat memiliki banyak alamat.
- Satu user dapat memiliki banyak keranjang.
- Satu user dapat membuat banyak pesanan.
- Satu kategori dapat memiliki banyak produk.
- Satu produk dapat memiliki banyak gambar.
- Satu keranjang dapat memiliki banyak item.
- Satu produk dapat berada di banyak item keranjang.
- Satu pesanan dapat memiliki banyak item.
- Satu produk dapat berada di banyak item pesanan.
- Satu pesanan dapat memiliki satu atau beberapa pembayaran.

---

## 3. ERD Diagram

```mermaid
erDiagram

    USERS ||--o{ ADDRESSES : memiliki
    USERS ||--o{ CARTS : memiliki
    USERS ||--o{ ORDERS : membuat

    CATEGORIES ||--o{ PRODUCTS : memiliki
    PRODUCTS ||--o{ PRODUCT_IMAGES : memiliki

    CARTS ||--o{ CART_ITEMS : berisi
    PRODUCTS ||--o{ CART_ITEMS : digunakan

    ORDERS ||--o{ ORDER_ITEMS : memiliki
    PRODUCTS ||--o{ ORDER_ITEMS : dibeli

    ORDERS ||--o{ PAYMENTS : memiliki

    USERS {
        uuid id PK
        varchar name
        varchar email UK
        varchar password_hash
        varchar role
        boolean is_active
        timestamp created_at
        timestamp updated_at
    }

    ADDRESSES {
        uuid id PK
        uuid user_id FK
        varchar label
        varchar recipient_name
        varchar phone
        text address
        varchar city
        varchar province
        varchar postal_code
        boolean is_default
        timestamp created_at
        timestamp updated_at
    }

    CATEGORIES {
        uuid id PK
        varchar name
        varchar slug UK
        text description
        boolean is_active
        timestamp created_at
        timestamp updated_at
    }

    PRODUCTS {
        uuid id PK
        uuid category_id FK
        varchar name
        varchar slug UK
        varchar sku UK
        text description
        numeric price
        integer stock
        integer weight
        boolean is_active
        timestamp created_at
        timestamp updated_at
    }

    PRODUCT_IMAGES {
        uuid id PK
        uuid product_id FK
        text image_url
        boolean is_primary
        integer sort_order
        timestamp created_at
    }

    CARTS {
        uuid id PK
        uuid user_id FK
        varchar status
        timestamp created_at
        timestamp updated_at
    }

    CART_ITEMS {
        uuid id PK
        uuid cart_id FK
        uuid product_id FK
        integer quantity
        numeric unit_price
        timestamp created_at
        timestamp updated_at
    }

    ORDERS {
        uuid id PK
        varchar order_number UK
        uuid user_id FK
        varchar status
        varchar payment_status
        numeric subtotal
        numeric shipping_cost
        numeric discount_amount
        numeric total_amount
        varchar shipping_name
        varchar shipping_phone
        text shipping_address
        timestamp created_at
        timestamp updated_at
    }

    ORDER_ITEMS {
        uuid id PK
        uuid order_id FK
        uuid product_id FK
        varchar product_name
        varchar sku
        numeric unit_price
        integer quantity
        numeric subtotal
    }

    PAYMENTS {
        uuid id PK
        uuid order_id FK
        varchar payment_method
        varchar transaction_id
        numeric amount
        varchar status
        timestamp paid_at
        timestamp expired_at
        timestamp created_at
        timestamp updated_at
    }
```
