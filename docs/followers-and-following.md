# VieTick Follow System - Tài liệu API

## Tổng quan

Hệ thống Follow cho phép người dùng theo dõi (follow) và bỏ theo dõi (unfollow) những người dùng khác trong cộng đồng VieTick. Tính năng này giúp tăng tính tương tác và kết nối giữa các thành viên.

---

## 1. Database Schema

### Follow Table
```sql
CREATE TABLE follows (
    id CHAR(36) PRIMARY KEY,
    follower_id CHAR(36) NOT NULL,    -- Người thực hiện follow
    following_id CHAR(36) NOT NULL,   -- Người được follow
    created_at TIMESTAMP NOT NULL,
    
    INDEX idx_follower (follower_id),
    INDEX idx_following (following_id),
    UNIQUE KEY unique_follow (follower_id, following_id),
    
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### User Model Updates
```go
type User struct {
    // ... existing fields ...
    
    // Follow relationships
    Followers  []Follow `gorm:"foreignKey:FollowingID;references:ID"` // Những người follow mình
    Following  []Follow `gorm:"foreignKey:FollowerID;references:ID"`  // Những người mình follow
}
```

---

## 2. API Endpoints

### 2.1 Follow User
**Follow một người dùng khác**

- **Endpoint:** `POST /follows/:id`
- **Method:** POST
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Path Parameter:** `:id` - ID của user muốn follow
- **Response:**
  ```json
  {
    "id": "uuid",
    "follower_id": "uuid",
    "following_id": "uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "follower": {
      "id": "uuid",
      "username": "user123",
      "email": "user@example.com",
      "point": 100
    },
    "following": {
      "id": "uuid",
      "username": "target_user",
      "email": "target@example.com",
      "point": 250
    }
  }
  ```
- **Error Cases:**
  - `400` - "cannot follow yourself"
  - `400` - "user to follow not found"
  - `400` - "already following this user"

### 2.2 Unfollow User
**Bỏ follow một người dùng**

- **Endpoint:** `DELETE /follows/:id`
- **Method:** DELETE
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Path Parameter:** `:id` - ID của user muốn unfollow
- **Response:**
  ```json
  {
    "message": "unfollowed successfully"
  }
  ```
- **Error Cases:**
  - `400` - "cannot unfollow yourself"
  - `400` - "not following this user"

### 2.3 Check Following Status
**Kiểm tra xem có đang follow user khác không**

- **Endpoint:** `GET /follows/:id/check`
- **Method:** GET
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Path Parameter:** `:id` - ID của user cần kiểm tra
- **Response:**
  ```json
  {
    "is_following": true
  }
  ```

### 2.4 Get User's Followers
**Lấy danh sách những người follow user**

- **Endpoint:** `GET /follows/:id/followers`
- **Method:** GET
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Path Parameter:** `:id` - ID của user
- **Query Parameters:**
  - `page` (optional, default: 1)
  - `limit` (optional, default: 20)
- **Response:**
  ```json
  {
    "data": [
      {
        "id": "uuid",
        "username": "follower1",
        "email": "follower1@example.com",
        "point": 150,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 25,
    "page": 1,
    "limit": 20
  }
  ```

### 2.5 Get User's Following
**Lấy danh sách những người mà user đang follow**

- **Endpoint:** `GET /follows/:id/following`
- **Method:** GET
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Path Parameter:** `:id` - ID của user
- **Query Parameters:**
  - `page` (optional, default: 1)
  - `limit` (optional, default: 20)
- **Response:**
  ```json
  {
    "data": [
      {
        "id": "uuid",
        "username": "following1",
        "email": "following1@example.com",
        "point": 300,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 15,
    "page": 1,
    "limit": 20
  }
  ```

### 2.6 Get User's Follow Stats
**Lấy thống kê follow của user**

- **Endpoint:** `GET /follows/:id/stats`
- **Method:** GET
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Path Parameter:** `:id` - ID của user
- **Response:**
  ```json
  {
    "followers_count": 25,
    "following_count": 15
  }
  ```

### 2.7 Get My Follow Stats
**Lấy thống kê follow của user hiện tại**

- **Endpoint:** `GET /me/follows/stats`
- **Method:** GET
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Response:**
  ```json
  {
    "followers_count": 10,
    "following_count": 8
  }
  ```

---

## 3. Business Logic

### 3.1 Follow Rules
- ✅ User có thể follow bất kỳ user nào khác
- ❌ User không thể follow chính mình
- ❌ Không thể follow cùng một user nhiều lần
- ✅ Có thể unfollow bất kỳ lúc nào

### 3.2 Data Integrity
- **Cascade Delete:** Khi user bị xóa, tất cả follow relationships sẽ bị xóa
- **Unique Constraint:** Mỗi cặp (follower_id, following_id) chỉ tồn tại một lần
- **Indexes:** Tối ưu performance cho các query follow

### 3.3 Performance Considerations
- **Pagination:** Tất cả danh sách followers/following đều có pagination
- **Database Indexes:** Index trên follower_id và following_id
- **Efficient Queries:** Sử dụng JOIN thay vì subqueries

---

## 4. Ví dụ sử dụng

### 4.1 Follow User
```bash
curl -X POST http://localhost:8080/follows/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### 4.2 Unfollow User
```bash
curl -X DELETE http://localhost:8080/follows/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### 4.3 Check Following Status
```bash
curl -X GET http://localhost:8080/follows/123e4567-e89b-12d3-a456-426614174000/check \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### 4.4 Get Followers
```bash
curl -X GET "http://localhost:8080/follows/123e4567-e89b-12d3-a456-426614174000/followers?page=1&limit=10" \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### 4.5 Get Following
```bash
curl -X GET "http://localhost:8080/follows/123e4567-e89b-12d3-a456-426614174000/following?page=1&limit=10" \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### 4.6 Get Follow Stats
```bash
curl -X GET http://localhost:8080/follows/123e4567-e89b-12d3-a456-426614174000/stats \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

---

## 5. Frontend Integration

### 5.1 Follow Button States
```javascript
// Check if following
const checkFollowing = async (userId) => {
  const response = await fetch(`/follows/${userId}/check`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  const data = await response.json();
  return data.is_following;
};

// Follow user
const followUser = async (userId) => {
  const response = await fetch(`/follows/${userId}`, {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${token}` }
  });
  return response.json();
};

// Unfollow user
const unfollowUser = async (userId) => {
  const response = await fetch(`/follows/${userId}`, {
    method: 'DELETE',
    headers: { 'Authorization': `Bearer ${token}` }
  });
  return response.json();
};
```

### 5.2 Follow Button Component
```jsx
const FollowButton = ({ userId, initialFollowing = false }) => {
  const [isFollowing, setIsFollowing] = useState(initialFollowing);
  const [isLoading, setIsLoading] = useState(false);

  const handleFollow = async () => {
    setIsLoading(true);
    try {
      if (isFollowing) {
        await unfollowUser(userId);
        setIsFollowing(false);
      } else {
        await followUser(userId);
        setIsFollowing(true);
      }
    } catch (error) {
      console.error('Follow action failed:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <button 
      onClick={handleFollow} 
      disabled={isLoading}
      className={isFollowing ? 'btn-unfollow' : 'btn-follow'}
    >
      {isLoading ? 'Loading...' : (isFollowing ? 'Unfollow' : 'Follow')}
    </button>
  );
};
```

---

## 6. Lưu ý quan trọng

### 6.1 Security
- ✅ Tất cả endpoints đều yêu cầu JWT authentication
- ✅ User chỉ có thể thao tác follow/unfollow của chính mình
- ✅ Validation đầy đủ cho user IDs

### 6.2 Performance
- ✅ Pagination cho tất cả danh sách
- ✅ Database indexes tối ưu
- ✅ Efficient queries với JOIN

### 6.3 Scalability
- ✅ Có thể mở rộng để thêm notifications
- ✅ Có thể thêm follow suggestions
- ✅ Có thể thêm follow analytics

### 6.4 Future Enhancements
- 🔮 Follow notifications
- 🔮 Follow suggestions based on interests
- 🔮 Follow analytics và insights
- 🔮 Follow groups hoặc topics
- 🔮 Follow privacy settings
