package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "vietick/internal/services"
)

type FollowController struct {
    followService *services.FollowService
}

func NewFollowController() *FollowController {
    return &FollowController{
        followService: services.NewFollowService(),
    }
}

// FollowUser thực hiện follow một user
func (c *FollowController) FollowUser(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    // Convert userID to uuid.UUID
    followerID, ok := userID.(uuid.UUID)
    if !ok {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
        return
    }

    // Get following user ID from URL parameter
    followingIDStr := ctx.Param("id")
    followingID, err := uuid.Parse(followingIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID to follow"})
        return
    }

    result, err := c.followService.FollowUser(followerID, followingID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, result)
}

// UnfollowUser thực hiện unfollow một user
func (c *FollowController) UnfollowUser(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    // Convert userID to uuid.UUID
    followerID, ok := userID.(uuid.UUID)
    if !ok {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
        return
    }

    // Get following user ID from URL parameter
    followingIDStr := ctx.Param("id")
    followingID, err := uuid.Parse(followingIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID to unfollow"})
        return
    }

    if err := c.followService.UnfollowUser(followerID, followingID); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "unfollowed successfully"})
}

// IsFollowing kiểm tra xem có đang follow user khác không
func (c *FollowController) IsFollowing(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    // Convert userID to uuid.UUID
    followerID, ok := userID.(uuid.UUID)
    if !ok {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
        return
    }

    // Get following user ID from URL parameter
    followingIDStr := ctx.Param("id")
    followingID, err := uuid.Parse(followingIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }

    isFollowing, err := c.followService.IsFollowing(followerID, followingID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"is_following": isFollowing})
}

// GetFollowers lấy danh sách followers của user
func (c *FollowController) GetFollowers(ctx *gin.Context) {
    userIDStr := ctx.Param("id")
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }

    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

    followers, total, err := c.followService.GetFollowers(userID, page, limit)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data":  followers,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

// GetFollowing lấy danh sách những người mà user đang follow
func (c *FollowController) GetFollowing(ctx *gin.Context) {
    userIDStr := ctx.Param("id")
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }

    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

    following, total, err := c.followService.GetFollowing(userID, page, limit)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data":  following,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

// GetUserFollowStats lấy thống kê follow của user
func (c *FollowController) GetUserFollowStats(ctx *gin.Context) {
    userIDStr := ctx.Param("id")
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }

    stats, err := c.followService.GetUserFollowStats(userID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, stats)
}

// GetMyFollowStats lấy thống kê follow của user hiện tại
func (c *FollowController) GetMyFollowStats(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    // Convert userID to uuid.UUID
    userIDUUID, ok := userID.(uuid.UUID)
    if !ok {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
        return
    }

    stats, err := c.followService.GetUserFollowStats(userIDUUID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, stats)
} 