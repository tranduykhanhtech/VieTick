# Vietick - Help Community Platform

## Quick Start

1. Clone repository:
```bash
git clone https://github.com/yourusername/vietick.git
cd vietick
```

2. Cài đặt dependencies:
```bash
go mod download
```

3. Cấu hình môi trường:
```bash
cp .env.example .env
# Chỉnh sửa các biến môi trường trong file .env
```

4. Chạy ứng dụng:
```bash
go run main.go
```

## API Endpoints

### Authentication
```bash
# Register ✅pass
curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d "{\"email\":\"user@example.com\",\"username\":\"username\",\"password\":\"password123\"}"

# Login ✅pass
curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d "{\"email\":\"user@example.com\",\"password\":\"password123\"}"
```

### Protected Routes (Cần token từ login)
```bash
# Get Profile ✅pass
curl -X GET http://localhost:8080/users/me -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYTQ3NjlmOWEtNDdmNy00ZWFlLTljMDctZDlkNDBjNTRiYzQ2IiwiZXhwIjoxNzQ5ODEwNDI1fQ.4jWJ2luU-h7QH6KBNzIs5C35Ru9MpfcFRpjjERlbRsU"

# Create Question ✅pass
curl -X POST http://localhost:8080/questions -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYTQ3NjlmOWEtNDdmNy00ZWFlLTljMDctZDlkNDBjNTRiYzQ2IiwiZXhwIjoxNzQ5ODEwNDI1fQ.4jWJ2luU-h7QH6KBNzIs5C35Ru9MpfcFRpjjERlbRsU" -H "Content-Type: application/json" -d "{\"title\":\"How to use Golang?\",\"content\":\"I am new to Golang...\"}"

# Get Questions ✅pass
curl -X GET "http://localhost:8080/questions?page=1&limit=10" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYTQ3NjlmOWEtNDdmNy00ZWFlLTljMDctZDlkNDBjNTRiYzQ2IiwiZXhwIjoxNzQ5ODEwNDI1fQ.4jWJ2luU-h7QH6KBNzIs5C35Ru9MpfcFRpjjERlbRsU"

# Create Answer ✅pass
curl -X POST http://localhost:8080/questions/84b22784-4b43-48f0-9206-a6eaedcc12fe/answers -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYTQ3NjlmOWEtNDdmNy00ZWFlLTljMDctZDlkNDBjNTRiYzQ2IiwiZXhwIjoxNzQ5ODEwNDI1fQ.4jWJ2luU-h7QH6KBNzIs5C35Ru9MpfcFRpjjERlbRsU" -H "Content-Type: application/json" -d "{\"content\":\"You should start with...\"}"

# Vote Answer
curl -X POST http://localhost:8080/answers/be9780b4-b1ae-4065-ab18-10cd881d5413/vote/up -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYTQ3NjlmOWEtNDdmNy00ZWFlLTljMDctZDlkNDBjNTRiYzQ2IiwiZXhwIjoxNzQ5ODExMDM5fQ.lQofFPkm84kdXoiVwOk-Jqz-teSqWb1LJFDnRKvAWvs"
```

## Cấu trúc Project
```
vietick/
├── config/         # Cấu hình ứng dụng
├── internal/       # Code nội bộ
│   ├── controllers/  # Xử lý HTTP requests
│   ├── middleware/   # Middleware (auth, logging, etc.)
│   ├── models/       # Database models
│   └── services/     # Business logic
├── pkg/           # Shared packages
├── routes/        # Route definitions
└── main.go        # Entry point
``` 