# Inventory Management API

## Cấu trúc thư mục
```
inventory-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── database.go
│   ├── models/
│   │   └── models.go
│   ├── handlers/
│   │   └── inventory.go
│   └── repository/
│       └── inventory_repo.go
├── docker-compose.yml
├── Dockerfile
├── .env.example
├── go.mod
└── README.md
```

## Hướng dẫn cài đặt

### 1. Khởi tạo project
```bash
mkdir inventory-api && cd inventory-api
go mod init inventory-api
```

### 2. Cài đặt dependencies
```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/danielgtaylor/huma/v2
go get github.com/danielgtaylor/huma/v2/adapters/humagin
go get github.com/joho/godotenv
```

### 3. Tạo file .env
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=inventory_db
SERVER_PORT=8080
```

### 4. Chạy ứng dụng

**Khởi động PostgreSQL với Docker:**
```bash
docker-compose up -d
```

**Chạy API:**
```bash
go run cmd/api/main.go
```

## API Endpoints

### Products
- `GET /api/products` - Lấy danh sách sản phẩm
- `GET /api/products/{id}` - Lấy chi tiết sản phẩm
- `POST /api/products` - Tạo sản phẩm mới
- `PUT /api/products/{id}` - Cập nhật sản phẩm
- `DELETE /api/products/{id}` - Xóa sản phẩm

### Inventory Transactions
- `GET /api/transactions` - Lấy danh sách giao dịch
- `POST /api/transactions` - Tạo giao dịch (nhập/xuất kho)

## Ví dụ Request

### Tạo sản phẩm mới
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Dell XPS 15",
    "sku": "DELL-XPS15-001",
    "description": "Laptop cao cấp",
    "price": 35000000,
    "quantity": 50
  }'
```

### Tạo giao dịch nhập kho
```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 10,
    "transaction_type": "IN",
    "notes": "Nhập hàng từ nhà cung cấp"
  }'
```

## Tính năng chính

- ✅ CRUD operations cho Products
- ✅ Quản lý giao dịch nhập/xuất kho
- ✅ Tự động cập nhật số lượng tồn kho
- ✅ API documentation với Huma
- ✅ Database migrations tự động
- ✅ Docker support
- ✅ Validation dữ liệu
- ✅ Error handling

## Technologies

- **Gin**: Web framework
- **GORM**: ORM
- **Huma**: API documentation & validation
- **PostgreSQL**: Database
- **Docker**: Containerization