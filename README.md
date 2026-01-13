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
JWT_SECRET=your-secret-key-change-in-production
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

### Authentication (Public)

- `POST /auth/register` - Đăng ký user mới
- `POST /auth/login` - Đăng nhập và nhận JWT token

### Users (Protected - Requires JWT)

- `GET /users/profile` - Xem profile của user hiện tại (any authenticated user)
- `POST /users/change-password` - Đổi mật khẩu (any authenticated user)
- `GET /users` - Lấy danh sách users (admin only)
- `GET /users/{id}` - Lấy thông tin user theo ID (admin or owner)
- `PUT /users/{id}` - Cập nhật user (admin or owner, only admin can change roles)
- `DELETE /users/{id}` - Xóa user (admin only, cannot delete self)

### Products

- `GET /products` - Lấy danh sách sản phẩm (public)
- `GET /products/{id}` - Lấy thông tin sản phẩm theo ID (public)
- `POST /products` - Tạo sản phẩm mới (authenticated users)
- `PUT /products/{id}` - Cập nhật sản phẩm (authenticated users)
- `DELETE /products/{id}` - Xóa sản phẩm (admin only)

### Transactions (Protected - Requires JWT)

- `POST /transactions` - Tạo giao dịch nhập/xuất kho
- `GET /transactions` - Lấy danh sách giao dịch (có phân trang)
- `GET /transactions/{id}` - Lấy thông tin giao dịch theo ID
- `GET /products/{id}/transactions` - Lấy lịch sử giao dịch của sản phẩm

## Ví dụ sử dụng

### Đăng ký user mới

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "password123",
    "email": "john@example.com",
    "phone": "0123456789",
    "role": "user"
  }'
```

**Password Requirements:**
- Minimum 8 characters
- Username: 3-50 characters, alphanumeric and underscore only

### Đăng nhập

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "password123"
  }'
```

Response sẽ trả về token:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "phone": "0123456789",
    "role": "user"
  }
}
```

### Xem profile (với JWT token)

```bash
curl -X GET http://localhost:8080/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

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

### Validation Rules

### Register User
- `username`: 3-50 ký tự, chỉ chữ cái, số và underscore
- `password`: tối thiểu 8 ký tự
- `email`: định dạng email hợp lệ
- `phone`: 10-15 ký tự, chỉ số và ký tự đặc biệt (+, -, space, ())
- `role`: chỉ nhận "admin" hoặc "user"

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

### Change Password
- `old_password`: bắt buộc
- `new_password`: tối thiểu 8 ký tự, phải khác mật khẩu cũ

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
- ✅ JWT Authentication với Bearer token
- ✅ Role-based access control (admin/user)
- ✅ Owner-based authorization
- ✅ Password hashing với bcrypt (minimum 8 characters)
- ✅ Secure error messages (no information leakage)
- ✅ Input validation với pattern matching
- ✅ Auto-generated OpenAPI documentation
- ✅ Soft delete cho products
- ✅ Transaction tracking (IN/OUT)
- ✅ Pagination support
- ✅ Docker support
- ✅ GORM ORM với PostgreSQL
- ✅ Clean Architecture (handler → service → repo)

## Security

Xem chi tiết về bảo mật tại [docs/SECURITY.md](docs/SECURITY.md)

### Quick Security Notes

- JWT tokens expire sau 24 giờ
- Passwords được hash với bcrypt
- Role-based access control (RBAC)
- Admin không thể tự xóa account của mình
- Users chỉ có thể update profile của mình (trừ role)
- Generic error messages để tránh user enumeration
- Products list/detail là public, các operations khác cần authentication

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
