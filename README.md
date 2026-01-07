# Inventory API

API quản lý kho hàng được xây dựng với Go, Huma framework, Gin và PostgreSQL.

## Yêu cầu

- Go 1.25+
- Docker & Docker Compose
- PostgreSQL 15 (nếu chạy local không dùng Docker)

## Cài đặt

### 1. Clone project và cài đặt dependencies

```bash
git clone <repository-url>
cd inventory-api
go mod download
```

### 2. Cấu hình môi trường

Tạo file `.env` từ file mẫu:

```bash
copy .env.example .env
```

Chỉnh sửa file `.env` theo môi trường của bạn:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=inventory_db
SERVER_PORT=8080
```

## Chạy ứng dụng

### Cách 1: Sử dụng Docker Compose (Khuyến nghị)

Chạy cả API và PostgreSQL trong Docker:

```bash
docker-compose up -d
```

API sẽ chạy tại: `http://localhost:8080`

Dừng services:

```bash
docker-compose down
```

Xóa cả volumes (database data):

```bash
docker-compose down -v
```

### Cách 2: Chạy local

#### Bước 1: Khởi động PostgreSQL

Sử dụng Docker:

```bash
docker-compose up -d postgres
```

Hoặc cài đặt PostgreSQL local và tạo database:

```sql
CREATE DATABASE inventory_db;
```

#### Bước 2: Chạy API

```bash
go run cmd/api/main.go
```

Hoặc build và chạy:

```bash
go build -o inventory-api.exe cmd/api/main.go
inventory-api.exe
```

## API Documentation

Sau khi chạy ứng dụng, truy cập:

- **API Docs (Interactive)**: http://localhost:8080/docs
- **OpenAPI JSON**: http://localhost:8080/openapi.json
- **OpenAPI YAML**: http://localhost:8080/openapi.yaml

## API Endpoints

### Products

- `POST /products` - Tạo sản phẩm mới
- `GET /products` - Lấy danh sách sản phẩm (có phân trang)
- `GET /products/{id}` - Lấy thông tin sản phẩm theo ID
- `PUT /products/{id}` - Cập nhật sản phẩm
- `DELETE /products/{id}` - Xóa sản phẩm (soft delete)

### Transactions

- `POST /transactions` - Tạo giao dịch nhập/xuất kho
- `GET /transactions` - Lấy danh sách giao dịch (có phân trang)
- `GET /transactions/{id}` - Lấy thông tin giao dịch theo ID
- `GET /products/{id}/transactions` - Lấy lịch sử giao dịch của sản phẩm

## Ví dụ sử dụng

### Tạo sản phẩm mới

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Dell XPS 15",
    "sku": "DELL-XPS-15",
    "description": "Laptop cao cấp",
    "price": 25000000,
    "quantity": 10
  }'
```

### Tạo giao dịch nhập kho

```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 5,
    "transaction_type": "IN",
    "notes": "Nhập hàng từ nhà cung cấp"
  }'
```

### Lấy danh sách sản phẩm

```bash
curl "http://localhost:8080/products?limit=10&offset=0"
```

## Validation Rules

### Create Product
- `name`: bắt buộc, 1-255 ký tự
- `sku`: bắt buộc, 1-100 ký tự, unique
- `price`: bắt buộc, phải >= 0.01 (không được = 0)
- `quantity`: bắt buộc, phải >= 1 (không được = 0)

### Update Product
- Tất cả fields đều optional
- `price`: có thể = 0
- `quantity`: có thể = 0

### Create Transaction
- `product_id`: bắt buộc
- `quantity`: bắt buộc, phải >= 1
- `transaction_type`: bắt buộc, chỉ nhận "IN" hoặc "OUT"

## Cấu trúc thư mục

```
inventory-api/
├── cmd/
│   └── api/
│       └── main.go           # Entry point
├── internal/
│   ├── config/              # Configuration
│   ├── database/            # Database connection
│   ├── handler/             # HTTP handlers
│   ├── models/              # Data models & DTOs
│   ├── repo/                # Repository layer
│   └── services/            # Business logic
├── docs/                    # Documentation
├── .env.example             # Environment variables template
├── docker-compose.yml       # Docker compose config
├── Dockerfile               # Docker build config
└── README.md
```

## Features

- ✅ RESTful API với Huma framework
- ✅ Auto-generated OpenAPI documentation
- ✅ Input validation tự động
- ✅ Soft delete cho products
- ✅ Transaction tracking (IN/OUT)
- ✅ Pagination support
- ✅ Docker support
- ✅ GORM ORM với PostgreSQL

## Troubleshooting

### Lỗi kết nối database

Kiểm tra:
1. PostgreSQL đã chạy chưa
2. Thông tin kết nối trong `.env` đúng chưa
3. Database đã được tạo chưa

### Port 8080 đã được sử dụng

Đổi `SERVER_PORT` trong file `.env` sang port khác (ví dụ: 8081)

## License

MIT
