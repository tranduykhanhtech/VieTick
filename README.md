# VieTick - Nền tảng Hỏi đáp Cộng đồng

VieTick là một nền tảng hỏi đáp cộng đồng được xây dựng bằng Go, cho phép người dùng đặt câu hỏi, trả lời, bình chọn và xác minh câu trả lời. Hệ thống được thiết kế với kiến trúc Clean Architecture, sử dụng JWT authentication và MySQL database.

## 🌟 Tính năng chính

### 👤 Quản lý người dùng
- Đăng ký và đăng nhập tài khoản
- Hệ thống điểm tích lũy
- Quản lý thông tin cá nhân
- JWT-based authentication

### ❓ Hệ thống hỏi đáp
- Tạo và quản lý câu hỏi
- Trả lời câu hỏi
- Xem danh sách câu hỏi và trả lời
- Phân trang và tìm kiếm

### 👍 Bình chọn và đánh giá
- Vote up/down cho câu trả lời
- Hệ thống điểm cho người dùng
- Thống kê số lượng vote

### ✅ Xác minh nội dung
- Xác minh câu trả lời
- Hệ thống báo cáo nội dung không phù hợp
- Quản lý chất lượng nội dung

## 🏗️ Kiến trúc hệ thống

```
vietick/
├── cmd/api/           # Entry point của ứng dụng
├── config/            # Cấu hình database và môi trường
├── internal/          # Code nội bộ (không export)
│   ├── controllers/   # Xử lý HTTP requests
│   ├── middleware/    # Middleware (auth, CORS, logging)
│   ├── models/        # Database models (GORM)
│   └── services/      # Business logic
├── pkg/               # Shared packages (JWT, logger, utils)
└── routes/            # Định nghĩa routes
```

## 🛠️ Công nghệ sử dụng

- **Backend**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **Database**: MySQL với GORM (ORM)
- **Authentication**: JWT (JSON Web Tokens)
- **Logging**: Zerolog
- **Environment**: Godotenv
- **UUID**: Google UUID library

## 📊 Mô hình dữ liệu

### User (Người dùng)
- ID, Email, Username, Password, Point
- Quan hệ: Questions, Answers, Votes

### Question (Câu hỏi)
- ID, Title, Content, UserID
- Quan hệ: User (người tạo), Answers

### Answer (Câu trả lời)
- ID, Content, QuestionID, UserID, IsVerified, VerifiedBy, Reported
- Quan hệ: Question, User (người trả lời), Verifier

### Vote (Bình chọn)
- ID, UserID, AnswerID, Type (up/down)
- Quan hệ: User, Answer

## 🚀 Quick Start

### 1. Clone repository
```bash
git clone https://github.com/yourusername/vietick.git
cd vietick
```

### 2. Cài đặt dependencies
```bash
go mod download
```

### 3. Cấu hình môi trường
```bash
cp .env.example .env
# Chỉnh sửa các biến môi trường trong file .env
```

**Các biến môi trường cần thiết:**
```env
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_password
DB_DATABASE=vietick
JWT_SECRET=your_jwt_secret
ENV=development
```

### 4. Chạy ứng dụng
```bash
# Development mode
make dev

# Hoặc chạy trực tiếp
go run cmd/api/main.go
```

## 📡 API Endpoints

### Authentication
```bash
# Đăng ký
curl -X POST https://vietick.onrender.com/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","username":"username","password":"password123"}'

# Đăng nhập
curl -X POST https://vietick.onrender.com/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Protected Routes (Cần JWT token)

#### 👤 User Management
```bash
# Lấy thông tin cá nhân
curl -X GET http://localhost:8080/users/me \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

#### ❓ Question Management
```bash
# Tạo câu hỏi
curl -X POST http://localhost:8080/questions \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"How to use Golang?","content":"I am new to Golang..."}'

# Lấy danh sách câu hỏi
curl -X GET "http://localhost:8080/questions?page=1&limit=10" \
  -H "Authorization: Bearer <JWT_TOKEN>"

# Lấy chi tiết câu hỏi
curl -X GET http://localhost:8080/questions/<question_id> \
  -H "Authorization: Bearer <JWT_TOKEN>"

# Cập nhật câu hỏi
curl -X PUT http://localhost:8080/questions/<question_id> \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated title","content":"Updated content"}'

# Xóa câu hỏi
curl -X DELETE http://localhost:8080/questions/<question_id> \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

#### 💬 Answer Management
```bash
# Tạo câu trả lời
curl -X POST http://localhost:8080/questions/<question_id>/answers \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"content":"This is my answer..."}'

# Lấy danh sách câu trả lời
curl -X GET "http://localhost:8080/questions/<question_id>/answers?page=1&limit=10" \
  -H "Authorization: Bearer <JWT_TOKEN>"

# Xác minh câu trả lời
curl -X POST http://localhost:8080/answers/<answer_id>/verify \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

#### 👍 Vote Management
```bash
# Vote up cho câu trả lời
curl -X POST http://localhost:8080/answers/<answer_id>/vote/up \
  -H "Authorization: Bearer <JWT_TOKEN>"

# Vote down cho câu trả lời
curl -X POST http://localhost:8080/answers/<answer_id>/vote/down \
  -H "Authorization: Bearer <JWT_TOKEN>"

# Lấy số lượng vote
curl -X GET http://localhost:8080/answers/<answer_id>/votes \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

## 🔧 Development Commands

```bash
# Build ứng dụng
make build

# Chạy tests
make test

# Chạy linter
make lint

# Clean build files
make clean

# Install dependencies
make deps

# Database migration
make migrate

# Development mode
make dev

# Production mode
make prod
```

## 🌐 Deployment

### Production
- **URL**: https://vietick.onrender.com
- **Platform**: Render.com
- **Database**: MySQL (managed)

### Local Development
- **Port**: 8080
- **Database**: MySQL local
- **Environment**: Development mode

## 🔐 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: Bcrypt password encryption
- **CORS Protection**: Cross-origin resource sharing middleware
- **Input Validation**: Request validation and sanitization
- **Error Handling**: Comprehensive error handling and logging

## 📈 Performance Features

- **Database Connection Pooling**: Optimized database connections
- **Pagination**: Efficient data loading with pagination
- **Indexing**: Database indexes for better query performance
- **Caching**: Ready for Redis integration

## 🧪 Testing

```bash
# Chạy tất cả tests
go test -v ./...

# Chạy tests với coverage
go test -v -cover ./...
```

## 📝 API Documentation

Xem chi tiết API documentation tại: [VieTick_API_Documentation.md](./VieTick_API_Documentation.md)

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

Nếu bạn gặp vấn đề hoặc có câu hỏi, vui lòng tạo issue trên GitHub repository.

---

**VieTick** - Nền tảng hỏi đáp cộng đồng mạnh mẽ và dễ sử dụng! 🚀 