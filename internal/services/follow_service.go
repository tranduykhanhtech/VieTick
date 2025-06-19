package services

import (
    "errors"
    "time"

    "github.com/google/uuid"
    "vietick/config"
    "vietick/internal/models"
)

type FollowService struct{}

type FollowResponse struct {
    ID          uuid.UUID `json:"id"`
    FollowerID  uuid.UUID `json:"follower_id"`
    FollowingID uuid.UUID `json:"following_id"`
    CreatedAt   time.Time `json:"created_at"`
    Follower    models.User `json:"follower"`
    Following   models.User `json:"following"`
}

type UserFollowStats struct {
    FollowersCount int64 `json:"followers_count"`
    FollowingCount int64 `json:"following_count"`
}

func NewFollowService() *FollowService {
    return &FollowService{}
}

// FollowUser thực hiện follow một user khác
func (s *FollowService) FollowUser(followerID, followingID uuid.UUID) (*FollowResponse, error) {
    // Không thể follow chính mình
    if followerID == followingID {
        return nil, errors.New("cannot follow yourself")
    }

    // Kiểm tra user được follow có tồn tại không
    var followingUser models.User
    if err := config.DB.First(&followingUser, "id = ?", followingID).Error; err != nil {
        return nil, errors.New("user to follow not found")
    }

    // Kiểm tra đã follow chưa
    var existingFollow models.Follow
    if err := config.DB.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&existingFollow).Error; err == nil {
        return nil, errors.New("already following this user")
    }

    // Tạo follow relationship
    now := time.Now()
    follow := models.Follow{
        FollowerID:  followerID,
        FollowingID: followingID,
        CreatedAt:   now,
    }

    if err := config.DB.Create(&follow).Error; err != nil {
        return nil, err
    }

    // Load follow với user data
    if err := config.DB.Preload("Follower").Preload("Following").First(&follow, follow.ID).Error; err != nil {
        return nil, err
    }

    // Gửi notification đến user được follow
    notificationService := NewNotificationService()
    notificationService.SendNotificationToUser(
        followingID,
        models.NotificationTypeFollow,
        "Bạn có người theo dõi mới!",
        follow.Follower.Username+" đã follow bạn.",
        map[string]interface{}{
            "follower_id":   follow.FollowerID,
            "follower_name": follow.Follower.Username,
        },
    )

    return &FollowResponse{
        ID:          follow.ID,
        FollowerID:  follow.FollowerID,
        FollowingID: follow.FollowingID,
        CreatedAt:   follow.CreatedAt,
        Follower:    follow.Follower,
        Following:   follow.Following,
    }, nil
}

// UnfollowUser thực hiện unfollow một user
func (s *FollowService) UnfollowUser(followerID, followingID uuid.UUID) error {
    // Không thể unfollow chính mình
    if followerID == followingID {
        return errors.New("cannot unfollow yourself")
    }

    // Kiểm tra follow relationship có tồn tại không
    var follow models.Follow
    if err := config.DB.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&follow).Error; err != nil {
        return errors.New("not following this user")
    }

    // Xóa follow relationship
    if err := config.DB.Delete(&follow).Error; err != nil {
        return err
    }

    // Gửi notification đến user bị unfollow (optional)
    notificationService := NewNotificationService()
    notificationService.SendNotificationToUser(
        followingID,
        models.NotificationTypeUnfollow,
        "Bạn vừa bị unfollow",
        follow.Follower.Username+" đã unfollow bạn.",
        map[string]interface{}{
            "follower_id":   follow.FollowerID,
            "follower_name": follow.Follower.Username,
        },
    )

    return nil
}

// IsFollowing kiểm tra xem một user có đang follow user khác không
func (s *FollowService) IsFollowing(followerID, followingID uuid.UUID) (bool, error) {
    var count int64
    if err := config.DB.Model(&models.Follow{}).
        Where("follower_id = ? AND following_id = ?", followerID, followingID).
        Count(&count).Error; err != nil {
        return false, err
    }
    return count > 0, nil
}

// GetFollowers lấy danh sách followers của một user
func (s *FollowService) GetFollowers(userID uuid.UUID, page, limit int) ([]models.User, int64, error) {
    var followers []models.User
    var total int64

    // Get total count
    if err := config.DB.Model(&models.Follow{}).
        Where("following_id = ?", userID).
        Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Get followers with pagination
    offset := (page - 1) * limit
    if err := config.DB.Model(&models.User{}).
        Joins("JOIN follows ON users.id = follows.follower_id").
        Where("follows.following_id = ?", userID).
        Order("follows.created_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&followers).Error; err != nil {
        return nil, 0, err
    }

    return followers, total, nil
}

// GetFollowing lấy danh sách những người mà user đang follow
func (s *FollowService) GetFollowing(userID uuid.UUID, page, limit int) ([]models.User, int64, error) {
    var following []models.User
    var total int64

    // Get total count
    if err := config.DB.Model(&models.Follow{}).
        Where("follower_id = ?", userID).
        Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Get following with pagination
    offset := (page - 1) * limit
    if err := config.DB.Model(&models.User{}).
        Joins("JOIN follows ON users.id = follows.following_id").
        Where("follows.follower_id = ?", userID).
        Order("follows.created_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&following).Error; err != nil {
        return nil, 0, err
    }

    return following, total, nil
}

// GetUserFollowStats lấy thống kê follow của user
func (s *FollowService) GetUserFollowStats(userID uuid.UUID) (*UserFollowStats, error) {
    var followersCount, followingCount int64

    // Count followers
    if err := config.DB.Model(&models.Follow{}).
        Where("following_id = ?", userID).
        Count(&followersCount).Error; err != nil {
        return nil, err
    }

    // Count following
    if err := config.DB.Model(&models.Follow{}).
        Where("follower_id = ?", userID).
        Count(&followingCount).Error; err != nil {
        return nil, err
    }

    return &UserFollowStats{
        FollowersCount: followersCount,
        FollowingCount: followingCount,
    }, nil
}

// GetMutualFollowers lấy danh sách mutual followers (cả hai follow nhau)
func (s *FollowService) GetMutualFollowers(userID1, userID2 uuid.UUID) ([]models.User, error) {
    var mutualUsers []models.User

    if err := config.DB.Model(&models.User{}).
        Joins("JOIN follows f1 ON users.id = f1.follower_id").
        Joins("JOIN follows f2 ON users.id = f2.following_id").
        Where("f1.following_id = ? AND f2.follower_id = ?", userID1, userID2).
        Find(&mutualUsers).Error; err != nil {
        return nil, err
    }

    return mutualUsers, nil
} 