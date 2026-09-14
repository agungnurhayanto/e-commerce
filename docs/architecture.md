# Application Architecture

## 1. Architecture Style

Aplikasi menggunakan Modular Monolith Architecture.

Backend dibuat dalam satu aplikasi Golang, tetapi setiap fitur dipisahkan berdasarkan modul.

Contoh modul:

- Auth
- User
- Category
- Product
- Cart
- Order
- Payment

## 2. Layer Architecture

Setiap modul menggunakan beberapa layer:

```text
Handler
   ↓
Service
   ↓
Repository
   ↓
Database
```
