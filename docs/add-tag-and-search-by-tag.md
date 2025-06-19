# VieTick API Documentation - Cập nhật Tags & Search

## Tổng quan

VieTick là hệ thống hỏi đáp cộng đồng với tính năng tags và tìm kiếm nâng cao. Tài liệu này mô tả các tính năng mới được thêm vào.

---

## 1. Tính năng Tags (Thẻ)

### 1.1 Tạo Tag mới
- **Endpoint:** `POST /tags`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Body:**
  ```json
  {
    "name": "golang",
    "description": "Ngôn ngữ lập trình Go",
    "color": "#00ADD8"
  }
  ```
- **Response:**
  ```json
  {
    "ID": "uuid",
    "Name": "golang",
    "Description": "Ngôn ngữ lập trình Go",
    "Color": "#00ADD8",
    "UsageCount": 0,
    "CreatedAt": "2024-01-01T00:00:00Z",
    "UpdatedAt": "2024-01-01T00:00:00Z"
  }
  ```

### 1.2 Lấy danh sách Tags
- **Endpoint:** `GET /tags`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Query:** `?page=1&limit=20`
- **Response:**
  ```json
  {
    "data": [
      {
        "ID": "uuid",
        "Name": "golang",
        "Description": "Ngôn ngữ lập trình Go",
        "Color": "#00ADD8",
        "UsageCount": 15,
        "CreatedAt": "2024-01-01T00:00:00Z",
        "UpdatedAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "limit": 20
  }
  ```

### 1.3 Lấy Tag theo ID
- **Endpoint:** `GET /tags/:id`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Response:** Trả về object tag

### 1.4 Cập nhật Tag
- **Endpoint:** `PUT /tags/:id`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Body:**
  ```json
  {
    "name": "go-language",
    "description": "Ngôn ngữ lập trình Go (cập nhật)",
    "color": "#00ADD8"
  }
  ```

### 1.5 Xóa Tag
- **Endpoint:** `DELETE /tags/:id`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Response:**
  ```json
  {
    "message": "tag deleted successfully"
  }
  ```

---

## 2. Tính năng Search (Tìm kiếm)

### 2.1 Tìm kiếm Tags
- **Endpoint:** `GET /search/tags`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Query:** `?q=golang&limit=10`
- **Response:**
  ```json
  {
    "data": [
      {
        "ID": "uuid",
        "Name": "golang",
        "Description": "Ngôn ngữ lập trình Go",
        "Color": "#00ADD8",
        "UsageCount": 15
      }
    ]
  }
  ```

### 2.2 Tìm kiếm Câu hỏi theo từ khóa
- **Endpoint:** `GET /search/questions`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Query:** `?q=how to use golang&page=1&limit=10`
- **Response:**
  ```json
  {
    "data": [
      {
        "ID": "uuid",
        "Title": "How to use Golang?",
        "Content": "I am new to Golang...",
        "UserID": "uuid",
        "Tags": [
          {
            "ID": "uuid",
            "Name": "golang",
            "Color": "#00ADD8"
          }
        ],
        "User": {
          "ID": "uuid",
          "Username": "user123"
        },
        "CreatedAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 25,
    "page": 1,
    "limit": 10,
    "query": "how to use golang"
  }
  ```

### 2.3 Tìm kiếm Câu hỏi theo Tag
- **Endpoint:** `GET /search/questions/tag/:tag`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Query:** `?page=1&limit=10`
- **Response:**
  ```json
  {
    "data": [
      {
        "ID": "uuid",
        "Title": "Golang best practices",
        "Content": "What are the best practices...",
        "UserID": "uuid",
        "Tags": [
          {
            "ID": "uuid",
            "Name": "golang",
            "Color": "#00ADD8"
          }
        ],
        "User": {
          "ID": "uuid",
          "Username": "user123"
        },
        "CreatedAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 15,
    "page": 1,
    "limit": 10,
    "tag": "golang"
  }
  ```

---

## 3. Cập nhật API Câu hỏi với Tags

### 3.1 Tạo câu hỏi với Tags
- **Endpoint:** `POST /questions`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Body:**
  ```json
  {
    "title": "How to use Golang?",
    "content": "I am new to Golang and want to learn...",
    "tags": ["golang", "programming", "beginner"]
  }
  ```
- **Response:** Trả về câu hỏi với tags đã được gắn

### 3.2 Cập nhật câu hỏi với Tags
- **Endpoint:** `PUT /questions/:id`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Body:**
  ```json
  {
    "title": "Updated title",
    "content": "Updated content",
    "tags": ["golang", "advanced", "concurrency"]
  }
  ```

### 3.3 Lấy câu hỏi với Tags
- **Endpoint:** `GET /questions` hoặc `GET /questions/:id`
- **Response:** Bao gồm thông tin tags trong response

---

## 4. Đặc điểm kỹ thuật

### 4.1 Tag Features
- **Normalization:** Tên tag được chuẩn hóa (lowercase, trim spaces)
- **Unique Names:** Tên tag phải duy nhất
- **Usage Count:** Tự động đếm số lần sử dụng
- **Color Support:** Hỗ trợ màu sắc cho UI
- **Auto Creation:** Tự động tạo tag mới nếu chưa tồn tại

### 4.2 Search Features
- **Fuzzy Search:** Tìm kiếm mờ cho tags
- **Full-text Search:** Tìm kiếm trong title và content
- **Tag-based Filtering:** Lọc câu hỏi theo tag
- **Pagination:** Hỗ trợ phân trang cho tất cả kết quả

### 4.3 Database Schema
```sql
-- Tags table
CREATE TABLE tags (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    color VARCHAR(7) DEFAULT '#007bff',
    usage_count BIGINT DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Question-Tags many-to-many relationship
CREATE TABLE question_tags (
    question_id CHAR(36),
    tag_id CHAR(36),
    PRIMARY KEY (question_id, tag_id),
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);
```

---

## 5. Ví dụ sử dụng

### 5.1 Tạo câu hỏi với tags
```bash
curl -X POST http://localhost:8080/questions \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Golang concurrency patterns",
    "content": "What are the best concurrency patterns in Go?",
    "tags": ["golang", "concurrency", "patterns"]
  }'
```

### 5.2 Tìm kiếm theo tag
```bash
curl -X GET "http://localhost:8080/search/questions/tag/golang?page=1&limit=10" \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### 5.3 Tìm kiếm tags
```bash
curl -X GET "http://localhost:8080/search/tags?q=go&limit=5" \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### 5.4 Tìm kiếm câu hỏi
```bash
curl -X GET "http://localhost:8080/search/questions?q=concurrency&page=1&limit=10" \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

---

## 6. Lưu ý quan trọng

1. **Tag Names:** Tên tag được chuẩn hóa tự động (lowercase, trim spaces)
2. **Usage Count:** Số lần sử dụng được cập nhật tự động khi tạo/xóa câu hỏi
3. **Transactions:** Tất cả operations với tags đều sử dụng database transactions
4. **Validation:** Tags có validation cho độ dài và format
5. **Performance:** Search được tối ưu với database indexes
6. **Security:** Tất cả endpoints đều yêu cầu JWT authentication
