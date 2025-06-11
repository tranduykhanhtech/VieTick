# VieTick API Documentation

## Tổng quan

VieTick là hệ thống hỏi đáp, hỗ trợ xác thực người dùng, đặt câu hỏi, trả lời, vote, xác minh câu trả lời, v.v.  
Tất cả các API trả về dữ liệu dạng JSON.

---

## 1. Đăng ký & Đăng nhập

```
URL: https://vietick.onrender.com
```

### Đăng ký tài khoản

- **Endpoint:** `POST /register`
- **Body:**
  ```json
  {
    "email": "user@example.com",
    "username": "username",
    "password": "password123"
  }
  ```
- **Response:**
  ```json
  {
    "token": "JWT_TOKEN",
    "user": { ... }
  }
  ```

### Đăng nhập

- **Endpoint:** `POST /login`
- **Body:**
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- **Response:**
  ```json
  {
    "token": "JWT_TOKEN",
    "user": { ... }
  }
  ```

---

## 2. Xác thực (Authentication)

- Tất cả các API cần xác thực phải gửi header:
  ```
  Authorization: Bearer <JWT_TOKEN>
  ```

---

## 3. Người dùng

### Lấy thông tin cá nhân

- **Endpoint:** `GET /users/me`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Response:**
  ```json
  {
    "ID": "...",
    "Email": "...",
    "Username": "...",
    "Point": 0,
    ...
  }
  ```

---

## 4. Câu hỏi (Questions)

### Tạo câu hỏi

- **Endpoint:** `POST /questions`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Body:**
  ```json
  {
    "title": "Tiêu đề câu hỏi",
    "content": "Nội dung chi tiết"
  }
  ```
- **Response:** Trả về object câu hỏi vừa tạo.

### Lấy danh sách câu hỏi

- **Endpoint:** `GET /questions`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Query:** `?page=1&limit=10` (tùy chọn)
- **Response:** Danh sách câu hỏi, phân trang.

### Lấy chi tiết câu hỏi

- **Endpoint:** `GET /questions/:id`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Response:** Chi tiết câu hỏi.

---

## 5. Trả lời (Answers)

### Tạo trả lời cho câu hỏi

- **Endpoint:** `POST /questions/:id/answers`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Body:**
  ```json
  {
    "content": "Nội dung trả lời"
  }
  ```
- **Response:** Trả về object câu trả lời vừa tạo.

### Lấy danh sách trả lời cho câu hỏi

- **Endpoint:** `GET /questions/:id/answers`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Query:** `?page=1&limit=10` (tùy chọn)
- **Response:** Danh sách trả lời, phân trang.

---

## 6. Vote cho trả lời

### Vote up/down cho trả lời

- **Endpoint:** `POST /answers/:id/vote/:type`
  - `:type` là `up` hoặc `down`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Response:** Trả về object vote hoặc `{ "message": "Vote removed" }` nếu bỏ vote.

### Lấy số lượng vote

- **Endpoint:** `GET /answers/:id/votes`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Response:**
  ```json
  {
    "up_votes": 5,
    "down_votes": 2
  }
  ```

---

## 7. Xác minh trả lời

### Xác minh hoặc bỏ xác minh trả lời (chỉ dành cho user đặc biệt hoặc tự động)

- **Endpoint:** `POST /answers/:id/verify`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Response:**
  ```json
  {
    "message": "Answer verified successfully"
  }
  ```

---

## 8. Các lưu ý khác

- **Tất cả các API trả về lỗi đều có dạng:**
  ```json
  { "error": "Thông báo lỗi" }
  ```
- **Phân trang:** Sử dụng query `?page=1&limit=10` cho các API trả về danh sách.
- **JWT Token:** Lưu vào localStorage/sessionStorage phía frontend để sử dụng cho các request tiếp theo.

---

## 9. Ví dụ sử dụng với curl

**Đăng ký:**
```bash
curl -X POST https://vietick.onrender.com/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","username":"username","password":"password123"}'
```

**Đăng nhập:**
```bash
curl -X POST https://vietick.onrender.com/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

**Tạo câu hỏi:**
```bash
curl -X POST https://vietick.onrender.com/questions \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Tiêu đề","content":"Nội dung"}'
```

---

Nếu cần bổ sung chi tiết về từng response, error, hoặc luồng xác thực, hãy chỉnh sửa file này hoặc liên hệ backend để được hỗ trợ thêm! 