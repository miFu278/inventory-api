# Troubleshooting Guide

## Vấn đề đã gặp và giải pháp

### 1. Port Conflict (Port 5432 đã được sử dụng)

**Lỗi:**
```
Error response from daemon: Bind for 0.0.0.0:5432 failed: port is already allocated
```

**Nguyên nhân:**
- Port 5432 đã được sử dụng bởi container PostgreSQL khác (mini-jira-db)
- Hoặc có process khác đang lắng nghe trên port 5432

**Giải pháp:**
Đổi port mapping trong docker-compose.yml và .env:

```yaml
# docker-compose.yml
services:
  postgres:
    ports:
      - "5433:5432"  # Đổi từ 5432:5432 sang 5433:5432
```

```env
# .env
DB_PORT=5433  # Đổi từ 5432 sang 5433
```

**Kiểm tra port đang sử dụng:**
```bash
# Windows
netstat -ano | findstr :5432

# Linux/Mac
lsof -i :5432
```

### 2. Huma API - Pointer Parameters Error

**Lỗi:**
```
panic: pointers are not supported for path/query/header parameters
```

**Nguyên nhân:**
Huma v2 không hỗ trợ pointer types cho path/query/header parameters. Trong code cũ:

```go
type ProductListQuery struct {
    SKU      *string  `query:"sku"`      // ❌ Pointer không được hỗ trợ
    Name     *string  `query:"name"`     // ❌ Pointer không được hỗ trợ
    MinPrice *float64 `query:"min_price"` // ❌ Pointer không được hỗ trợ
    MaxPrice *float64 `query:"max_price"` // ❌ Pointer không được hỗ trợ
}
```

**Giải pháp:**
Sử dụng value types và convert sang pointer khi cần:

```go
type ProductListQuery struct {
    SKU      string  `query:"sku"`      // ✅ Value type
    Name     string  `query:"name"`     // ✅ Value type
    MinPrice float64 `query:"min_price"` // ✅ Value type
    MaxPrice float64 `query:"max_price"` // ✅ Value type
}

func (q *ProductListQuery) ToProductFilter() *ProductFilter {
    filter := &ProductFilter{}
    
    // Chỉ set pointer nếu value không empty/zero
    if q.SKU != "" {
        filter.SKU = &q.SKU
    }
    if q.Name != "" {
        filter.Name = &q.Name
    }
    if q.MinPrice > 0 {
        filter.MinPrice = &q.MinPrice
    }
    if q.MaxPrice > 0 {
        filter.MaxPrice = &q.MaxPrice
    }
    
    return filter
}
```

**Lưu ý:**
- Path parameters: Không dùng pointer
- Query parameters: Không dùng pointer
- Header parameters: Không dùng pointer
- Request body: Có thể dùng pointer

### 3. Container Restart Loop

**Triệu chứng:**
```bash
docker ps
# STATUS: Restarting (2) 1 second ago
```

**Cách debug:**
```bash
# Xem logs của container
docker logs <container_name> --tail 50

# Xem logs realtime
docker logs -f <container_name>
```

**Nguyên nhân thường gặp:**
1. Application panic/crash
2. Database connection failed
3. Port conflict
4. Missing environment variables
5. Code compilation errors

## Các lệnh hữu ích

### Docker Commands

```bash
# Xem tất cả containers
docker ps -a

# Xem logs
docker logs <container_name>
docker logs -f <container_name>  # Follow logs

# Stop và remove containers
docker-compose down

# Start containers
docker-compose up -d

# Rebuild và start
docker-compose up -d --build

# Xem port mapping
docker ps --format "table {{.Names}}\t{{.Ports}}"

# Remove all stopped containers
docker container prune
```

### Database Commands

```bash
# Connect to PostgreSQL container
docker exec -it inventory_postgres psql -U postgres -d inventory_db

# Backup database
docker exec inventory_postgres pg_dump -U postgres inventory_db > backup.sql

# Restore database
docker exec -i inventory_postgres psql -U postgres inventory_db < backup.sql
```

### Go Commands

```bash
# Build
go build -o inventory-api.exe ./cmd/api

# Run
./inventory-api.exe

# Test
go test ./...

# Check for errors
go vet ./...

# Format code
go fmt ./...

# Download dependencies
go mod download

# Tidy dependencies
go mod tidy
```

## Kiểm tra API hoạt động

### 1. Health Check
```bash
curl http://localhost:8080/docs
```

### 2. List Products
```bash
curl http://localhost:8080/products
```

### 3. Create Product
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "sku": "TEST-001",
    "description": "Test description",
    "price": 99.99,
    "quantity": 10
  }'
```

### 4. Filter Products
```bash
# By name
curl "http://localhost:8080/products?name=laptop"

# By price range
curl "http://localhost:8080/products?min_price=100&max_price=500"

# By SKU
curl "http://localhost:8080/products?sku=TEST-001"
```

## Cấu hình hiện tại

### Ports
- API: `8080`
- PostgreSQL: `5433` (mapped to container's 5432)

### Database
- Host: `localhost` (khi run local) hoặc `postgres` (khi run trong Docker)
- Port: `5433`
- User: `postgres`
- Password: `postgres`
- Database: `inventory_db`

### API Documentation
- URL: http://localhost:8080/docs
- OpenAPI JSON: http://localhost:8080/openapi.json
- OpenAPI YAML: http://localhost:8080/openapi.yaml

## Best Practices

1. **Luôn check logs khi có lỗi:**
   ```bash
   docker logs inventory_api --tail 50
   ```

2. **Kiểm tra database connection:**
   ```bash
   docker exec -it inventory_postgres psql -U postgres -d inventory_db -c "SELECT 1"
   ```

3. **Rebuild khi thay đổi code:**
   ```bash
   docker-compose up -d --build
   ```

4. **Clean up khi cần:**
   ```bash
   docker-compose down -v  # Remove volumes
   docker system prune     # Clean up unused resources
   ```

5. **Backup data trước khi thay đổi schema:**
   ```bash
   docker exec inventory_postgres pg_dump -U postgres inventory_db > backup_$(date +%Y%m%d).sql
   ```
