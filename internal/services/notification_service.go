package services

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"vietick/config"
	"vietick/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationService struct {
	clients map[uuid.UUID]*Client
	mutex   sync.RWMutex
}

type Client struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Channel  chan []byte
	IsActive bool
}

type NotificationData struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Data      string    `json:"data,omitempty"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		clients: make(map[uuid.UUID]*Client),
	}
}

// AddClient thêm client vào SSE connection
func (s *NotificationService) AddClient(userID uuid.UUID) *Client {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	clientID := uuid.New()
	client := &Client{
		ID:       clientID,
		UserID:   userID,
		Channel:  make(chan []byte, 100),
		IsActive: true,
	}

	s.clients[clientID] = client
	log.Printf("Client added: %s for user: %s", clientID, userID)
	return client
}

// RemoveClient xóa client khỏi SSE connection
func (s *NotificationService) RemoveClient(clientID uuid.UUID) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if client, exists := s.clients[clientID]; exists {
		close(client.Channel)
		delete(s.clients, clientID)
		log.Printf("Client removed: %s", clientID)
	}
}

// SendNotificationToUser gửi notification đến một user cụ thể
func (s *NotificationService) SendNotificationToUser(userID uuid.UUID, notificationType models.NotificationType, title, message string, data map[string]interface{}) error {
	// Lưu notification vào database
	notification := models.Notification{
		UserID:    userID,
		Type:      notificationType,
		Title:     title,
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if data != nil {
		dataJSON, err := json.Marshal(data)
		if err != nil {
			return err
		}
		notification.Data = string(dataJSON)
	}

	if err := config.DB.Create(&notification).Error; err != nil {
		return err
	}

	// Gửi notification qua SSE
	s.sendToUser(userID, notification)
	return nil
}

// SendNotificationToFollowers gửi notification đến tất cả followers
func (s *NotificationService) SendNotificationToFollowers(userID uuid.UUID, notificationType models.NotificationType, title, message string, data map[string]interface{}) error {
	// Lấy danh sách followers
	var followers []models.User
	if err := config.DB.Model(&models.User{}).
		Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", userID).
		Find(&followers).Error; err != nil {
		return err
	}

	// Gửi notification đến từng follower
	for _, follower := range followers {
		s.SendNotificationToUser(follower.ID, notificationType, title, message, data)
	}

	return nil
}

// sendToUser gửi notification qua SSE đến user
func (s *NotificationService) sendToUser(userID uuid.UUID, notification models.Notification) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	notificationData := NotificationData{
		ID:        notification.ID,
		Type:      string(notification.Type),
		Title:     notification.Title,
		Message:   notification.Message,
		Data:      notification.Data,
		IsRead:    notification.IsRead,
		CreatedAt: notification.CreatedAt,
	}

	data, err := json.Marshal(notificationData)
	if err != nil {
		log.Printf("Error marshaling notification: %v", err)
		return
	}

	// Gửi đến tất cả clients của user
	for _, client := range s.clients {
		if client.UserID == userID && client.IsActive {
			select {
			case client.Channel <- data:
				// Successfully sent
			default:
				// Channel is full, mark client as inactive
				client.IsActive = false
			}
		}
	}
}

// SSEHandler xử lý SSE connection
func (s *NotificationService) SSEHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(400, gin.H{"error": "invalid user ID"})
		return
	}

	// Set headers cho SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	// Tạo client mới
	client := s.AddClient(userIDUUID)
	defer s.RemoveClient(client.ID)

	// Gửi initial connection message
	c.SSEvent("connected", gin.H{"message": "Connected to notifications"})

	// Listen for notifications
	for {
		select {
		case data := <-client.Channel:
			c.SSEvent("notification", string(data))
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			log.Printf("Client disconnected: %s", client.ID)
			return
		}
	}
}

// GetUserNotifications lấy danh sách notifications của user
func (s *NotificationService) GetUserNotifications(userID uuid.UUID, page, limit int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	// Get total count
	if err := config.DB.Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get notifications with pagination
	offset := (page - 1) * limit
	if err := config.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// MarkNotificationAsRead đánh dấu notification đã đọc
func (s *NotificationService) MarkNotificationAsRead(notificationID, userID uuid.UUID) error {
	return config.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error
}

// MarkAllNotificationsAsRead đánh dấu tất cả notifications đã đọc
func (s *NotificationService) MarkAllNotificationsAsRead(userID uuid.UUID) error {
	return config.DB.Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error
}

// GetUnreadCount lấy số lượng notifications chưa đọc
func (s *NotificationService) GetUnreadCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// DeleteNotification xóa notification
func (s *NotificationService) DeleteNotification(notificationID, userID uuid.UUID) error {
	return config.DB.Where("id = ? AND user_id = ?", notificationID, userID).
		Delete(&models.Notification{}).Error
}
