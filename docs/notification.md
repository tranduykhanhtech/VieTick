# VieTick Notification System - Tài liệu API

## 1. Tổng quan

Hệ thống notification của VieTick sử dụng Server-Sent Events (SSE) để gửi thông báo real-time đến người dùng khi có các sự kiện như: follow, unfollow, câu hỏi mới, trả lời mới, v.v. Notification được lưu vào database và gửi đến client ngay lập tức nếu đang online.

---

## 2. Database Schema

### Notification Table
```sql
CREATE TABLE notifications (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL, -- Người nhận
    type VARCHAR(20) NOT NULL, -- Loại thông báo
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSON,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

---

## 3. API Endpoints

### 3.1 SSE Stream
- **Endpoint:** `GET /notifications/stream`
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Kết nối nhận notification real-time qua SSE.
- **Client Example:**
  ```javascript
  const eventSource = new EventSource('/notifications/stream');
  eventSource.addEventListener('notification', (event) => {
    const notification = JSON.parse(event.data);
    // Hiển thị notification desktop
    if (Notification.permission === 'granted') {
      new Notification(notification.title, {
        body: notification.message,
        icon: '/icon.png'
      });
    }
    // Update UI
    updateNotificationBadge();
    addToNotificationList(notification);
  });
  ```

### 3.2 Lấy danh sách notification
- **Endpoint:** `GET /notifications`
- **Query:** `?page=1&limit=20`
- **Response:**
  ```json
  {
    "data": [
      {
        "id": "uuid",
        "type": "question",
        "title": "Câu hỏi mới từ người bạn follow",
        "message": "user123 vừa đăng câu hỏi: ...",
        "data": "{...}",
        "is_read": false,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "limit": 20
  }
  ```

### 3.3 Đánh dấu đã đọc
- **Endpoint:** `POST /notifications/:id/read`
- **Method:** POST
- **Description:** Đánh dấu một notification đã đọc

### 3.4 Đánh dấu tất cả đã đọc
- **Endpoint:** `POST /notifications/read-all`
- **Method:** POST
- **Description:** Đánh dấu tất cả notification đã đọc

### 3.5 Đếm số notification chưa đọc
- **Endpoint:** `GET /notifications/unread-count`
- **Response:**
  ```json
  { "unread_count": 3 }
  ```

### 3.6 Xóa notification
- **Endpoint:** `DELETE /notifications/:id`

---

## 4. Notification Types
- `follow` - Có người follow bạn
- `unfollow` - Có người unfollow bạn
- `question` - Người bạn follow đăng câu hỏi mới
- `answer` - Câu hỏi của bạn có trả lời mới
- `vote` - Câu trả lời/câu hỏi của bạn được vote
- `verify` - Câu trả lời được xác minh
- `tag` - Có câu hỏi mới với tag bạn quan tâm

---

## 5. Luồng hoạt động
1. User đăng nhập, frontend kết nối SSE.
2. Khi có sự kiện (follow, question, answer, ...), backend lưu notification vào DB và gửi qua SSE nếu user đang online.
3. Frontend nhận notification, hiển thị desktop notification và cập nhật UI.
4. User có thể xem danh sách, đánh dấu đã đọc, xóa notification.

---

## 6. Frontend Integration
- **Yêu cầu quyền notification:**
  ```javascript
  if (Notification.permission === 'default') {
    Notification.requestPermission();
  }
  ```
- **Hiển thị notification:**
  ```javascript
  new Notification('VieTick', { body: notification.message, icon: '/icon.png' });
  ```
- **Cập nhật badge, danh sách notification:**
  - Tăng/giảm số lượng chưa đọc
  - Thêm notification vào danh sách

---

## 7. Lưu ý
- SSE chỉ gửi notification khi user đang online, offline sẽ nhận qua API khi reload.
- Notification được lưu DB để không bị mất khi offline.
- Có thể mở rộng cho các loại notification khác.
- Tất cả endpoints đều yêu cầu JWT authentication.
