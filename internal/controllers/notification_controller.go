package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "vietick/internal/services"
)

type NotificationController struct {
    notificationService *services.NotificationService
}

func NewNotificationController() *NotificationController {
    return &NotificationController{
        notificationService: services.NewNotificationService(),
    }
}

// SSEStream xử lý SSE connection cho real-time notifications
func (c *NotificationController) SSEStream(ctx *gin.Context) {
    c.notificationService.SSEHandler(ctx)
}

// GetNotifications lấy danh sách notifications của user
func (c *NotificationController) GetNotifications(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    userIDUUID, ok := userID.(uuid.UUID)
    if !ok {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
        return
    }

    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

    notifications, total, err := c.notificationService.GetUserNotifications(userIDUUID, page, limit)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data":  notifications,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

// MarkAsRead đánh dấu notification đã đọc
func (c *NotificationController) MarkAsRead(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    userIDUUID, ok := userID.(uuid.UUID)
    if !ok {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
        return
    }

    notificationIDStr := ctx.Param("id")
    notificationID, err := uuid.Parse(notificationIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
        return
    }

    if err := c.notificationService.MarkNotificationAsRead(notificationID, userIDUUID); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}

// MarkAllAsRead đánh dấu tất cả notifications đã đọc
func (c *NotificationController) MarkAllAsRead(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    userIDUUID, ok := userID.(uuid.UUID)
    if !ok {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
        return
    }

    if err := c.notificationService.MarkAllNotificationsAsRead(userIDUUID); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "all notifications marked as read"})
}

// GetUnreadCount lấy số lượng notifications chưa đọc
func (c *NotificationController) GetUnreadCount(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    userIDUUID, ok := userID.(uuid.UUID)
    if !ok {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
        return
    }

    count, err := c.notificationService.GetUnreadCount(userIDUUID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"unread_count": count})
}

// DeleteNotification xóa notification
func (c *NotificationController) DeleteNotification(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    userIDUUID, ok := userID.(uuid.UUID)
    if !ok {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
        return
    }

    notificationIDStr := ctx.Param("id")
    notificationID, err := uuid.Parse(notificationIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
        return
    }

    if err := c.notificationService.DeleteNotification(notificationID, userIDUUID); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "notification deleted"})
} 