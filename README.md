- **Frontend** : Next.js 16 - port 3000
- **API Gateway** : Go Fiber - port 8000 [[main.go](./backend/api-gateway/main.go)]
- **Product Service** : Go Fiber - port 8001 [[main.go](./backend/product-service/main.go)]
- **Order Service** : Go Fiber - port 8002 [[main.go](./backend/order-service/main.go)]
- **Database** : MongoDB (2 collection terpisah: `products` & `orders`) bisa dilihat [[disini(order)](./backend/order-service/repository/order_repository.go)] dan [[disini(products)](./backend/product-service/repository/product_repository.go)]

## Tech Stack

| Komponen        | Technology                                                                                                                                          |
| --------------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| Backend Product | Go 1.25, [Fiber v2](https://pkg.go.dev/github.com/gofiber/fiber/v2), [MongoDB Driver](https://www.mongodb.com/docs/drivers/go/current/get-started/) |
| Backend Order   | Go 1.25, [Fiber v2](https://pkg.go.dev/github.com/gofiber/fiber/v2), [MongoDB Driver](https://www.mongodb.com/docs/drivers/go/current/get-started/) |
| Frontend        | Next.js 15, React 19, TypeScript, Tailwind CSS, shadcn/ui                                                                                           |
| Database        | MongoDB (local / Docker)                                                                                                                            |

## Prasyarat

- **Go** 1.25
- **Node.js** 24.11.0
- **Docker** 29.5.2
- **MongoDB** berjalan di local (atau via Docker)

## Cara Menjalankan

### 1. Clone Repository

```bash
git clone https://github.com/Yafiakmal/bharta-technical-test.git
cd bharta-technical-test

```

### 2. jalankan mongodb dengan docker

```bash
docker compose up -d
# atau
docker volume create mongo_data

docker run -d \
  --name mongodb \
  -p 27017:27017 \
  -v mongo_data:/data/db \
  mongo:7
```

### 3. jalankan backend

> buka terminal baru untuk menjalankan product service

```bash
cd backend/product-service/
# lalu install dependencies dan jalankan
go mod tidy
go run main.go
```

> buka terminal baru lagi untuk menjalankan order service

```bash
cd backend/order-service/
# lalu install dependencies dan jalankan
go mod tidy
go run main.go
```

> buka terminal baru lagi untuk menjalankan gateway

```bash
cd backend/api-gateway/
# lalu install dependencies dan jalankan
go mod tidy
go run main.go
```

### 4. jalankan frontend

```bash
cd frontend
npm install
npm run dev
```

> **.env** tidak perlu saya kira, biar hardcoded aja

## Daftar End Point

### 1. product service

| Method | Endpoint            | Description       | Contoh Body                                |
| ------ | ------------------- | ----------------- | ------------------------------------------ |
| POST   | /products           | Create product    | {"name":"Mouse","price":150000,"stock":10} |
| GET    | /products           | Get all products  | -                                          |
| GET    | /products/:id       | Get product by ID | -                                          |
| PUT    | /products/:id       | Update product    | {"name":"Mouse Wireless","price":175000}   |
| DELETE | /products/:id       | Delete product -  |
| PATCH  | /products/:id/stock | Update stock only | {"stock":5}                                |

### 2. order service

| Method | Endpoint                 | Description                       | Contoh Body                       |
| ------ | ------------------------ | --------------------------------- | --------------------------------- |
| POST   | /orders                  | Create order                      | {"productId":"...", "quantity":2} |
| GET    | /orders                  | Get all orders                    | -                                 |
| GET    | /orders/:id              | Get order by ID                   | -                                 |
| PUT    | /orders/:id              | Update order quantity             | {"quantity":5}                    |
| DELETE | /orders/:id              | Delete order (with stock restore) | -                                 |
| DELETE | /orders/force/:id        | Force delete (testing only)       | -                                 |
| DELETE | /orders/admin/delete-all | Delete all orders (testing only)  | -                                 |

### contoh curl product

```bash
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Wireless Mouse","price":250000,"stock":20}'
```

```bash
curl -X GET http://localhost:8000/products
```

```bash
curl -X GET http://localhost:8000/products/:id
```

```bash
curl -X PUT http://localhost:8000/products/:id \
  -H "Content-Type: application/json" \
  -d '{"price":275000}'
```

```bash
curl -X PATCH http://localhost:8000/products/:id/stock \
  -H "Content-Type: application/json" \
  -d '{"stock":15}'
```

```bash
curl -X DELETE http://localhost:8000/products/:id
```

### contoh curl order

```bash
curl -X POST http://localhost:8000/orders \
  -H "Content-Type: application/json" \
  -d '{"productId":":id","quantity":2}'
```

```bash
curl -X GET http://localhost:8000/orders
```

```bash
curl -X GET http://localhost:8000/orders/:id
```

```bash
curl -X PUT http://localhost:8000/orders/:id \
  -H "Content-Type: application/json" \
  -d '{"quantity":3}'
```

```bash
curl -X DELETE http://localhost:8000/orders/:id
```

```bash
# untuk kebutuhan testing)
curl -X DELETE http://localhost:8000/orders/force/:id \
  -H "X-Admin-Key: test-secret-key"
```

## Jawaban Pertanyaan Singkat

1. Kenapa kamu pilih desain collection MongoDB seperti itu? (embed atau reference, alasannya)

   > karena saya tidak terbiasa dengan nosql dan saya lebih sering pakai foreign key pada database sql. dan juga kalau saya pelajari sedikit behaviornya sepertinya reference bisa lebih konsisten ketika ada perubahan pada dokumen product.

2. Apa kelemahan cara order memanggil product service di solusimu? Kalau punya waktu lebih,kamu perbaiki bagaimana?
   > tidak atomic, saya kurang tahu di mongodb + golang bisa atau ngga. tapi dengan PrismaJS + postgres saya bisa atomic dengan $transaction. jujur saya kurang tau best practice dalam arsitektur microservice dengan api gateway. kalo saya punya waktu lebih mungkin saya akan lebih banyak riset.
3. Kalau jumlah service ini bertambah banyak ke depan, apa yang akan kamu siapkan di level arsitektur agar tetap rapi dan mudah dikelola?
   > yang pasti adalah struktur folder yang lebih tertata, kemudian kontrak program juga lebih terpusat, utils buat logging sama format response api yang konsisten juga perlu.
