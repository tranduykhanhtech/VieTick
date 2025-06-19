package services

import (
    "errors"
    "strings"
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
    "vietick/config"
    "vietick/internal/models"
)

type TagService struct{}

type CreateTagRequest struct {
    Name        string `json:"name" binding:"required,min=2,max=50"`
    Description string `json:"description"`
    Color       string `json:"color"`
}

type UpdateTagRequest struct {
    Name        string `json:"name" binding:"required,min=2,max=50"`
    Description string `json:"description"`
    Color       string `json:"color"`
}

func NewTagService() *TagService {
    return &TagService{}
}

// CreateTag tạo tag mới
func (s *TagService) CreateTag(req CreateTagRequest) (*models.Tag, error) {
    // Normalize tag name (lowercase, trim spaces)
    normalizedName := strings.ToLower(strings.TrimSpace(req.Name))
    
    // Check if tag already exists
    var existingTag models.Tag
    if err := config.DB.Where("name = ?", normalizedName).First(&existingTag).Error; err == nil {
        return nil, errors.New("tag already exists")
    }

    now := time.Now()
    tag := models.Tag{
        Name:        normalizedName,
        Description: req.Description,
        Color:       req.Color,
        UsageCount:  0,
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    if tag.Color == "" {
        tag.Color = "#007bff" // Default blue color
    }

    if err := config.DB.Create(&tag).Error; err != nil {
        return nil, err
    }

    return &tag, nil
}

// GetTags lấy danh sách tags với pagination
func (s *TagService) GetTags(page, limit int) ([]models.Tag, int64, error) {
    var tags []models.Tag
    var total int64

    // Get total count
    if err := config.DB.Model(&models.Tag{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Get tags with pagination, ordered by usage count
    offset := (page - 1) * limit
    if err := config.DB.Order("usage_count DESC, name ASC").
        Offset(offset).
        Limit(limit).
        Find(&tags).Error; err != nil {
        return nil, 0, err
    }

    return tags, total, nil
}

// GetTagByID lấy tag theo ID
func (s *TagService) GetTagByID(tagID uuid.UUID) (*models.Tag, error) {
    var tag models.Tag
    if err := config.DB.First(&tag, "id = ?", tagID).Error; err != nil {
        return nil, errors.New("tag not found")
    }
    return &tag, nil
}

// GetTagByName lấy tag theo tên
func (s *TagService) GetTagByName(name string) (*models.Tag, error) {
    var tag models.Tag
    normalizedName := strings.ToLower(strings.TrimSpace(name))
    if err := config.DB.Where("name = ?", normalizedName).First(&tag).Error; err != nil {
        return nil, errors.New("tag not found")
    }
    return &tag, nil
}

// UpdateTag cập nhật tag
func (s *TagService) UpdateTag(tagID uuid.UUID, req UpdateTagRequest) (*models.Tag, error) {
    var tag models.Tag
    if err := config.DB.First(&tag, "id = ?", tagID).Error; err != nil {
        return nil, errors.New("tag not found")
    }

    normalizedName := strings.ToLower(strings.TrimSpace(req.Name))
    
    // Check if new name conflicts with existing tag
    if normalizedName != tag.Name {
        var existingTag models.Tag
        if err := config.DB.Where("name = ? AND id != ?", normalizedName, tagID).First(&existingTag).Error; err == nil {
            return nil, errors.New("tag name already exists")
        }
    }

    tag.Name = normalizedName
    tag.Description = req.Description
    if req.Color != "" {
        tag.Color = req.Color
    }
    tag.UpdatedAt = time.Now()

    if err := config.DB.Save(&tag).Error; err != nil {
        return nil, err
    }

    return &tag, nil
}

// DeleteTag xóa tag
func (s *TagService) DeleteTag(tagID uuid.UUID) error {
    var tag models.Tag
    if err := config.DB.First(&tag, "id = ?", tagID).Error; err != nil {
        return errors.New("tag not found")
    }

    // Check if tag is being used
    if tag.UsageCount > 0 {
        return errors.New("cannot delete tag that is being used")
    }

    if err := config.DB.Delete(&tag).Error; err != nil {
        return err
    }

    return nil
}

// GetOrCreateTags tạo tags nếu chưa tồn tại, trả về danh sách tag IDs
func (s *TagService) GetOrCreateTags(tagNames []string) ([]uuid.UUID, error) {
    var tagIDs []uuid.UUID

    for _, name := range tagNames {
        normalizedName := strings.ToLower(strings.TrimSpace(name))
        if normalizedName == "" {
            continue
        }

        // Try to find existing tag
        var tag models.Tag
        if err := config.DB.Where("name = ?", normalizedName).First(&tag).Error; err != nil {
            // Tag doesn't exist, create new one
            now := time.Now()
            newTag := models.Tag{
                Name:        normalizedName,
                Description: "",
                Color:       "#007bff",
                UsageCount:  0,
                CreatedAt:   now,
                UpdatedAt:   now,
            }

            if err := config.DB.Create(&newTag).Error; err != nil {
                return nil, err
            }
            tagIDs = append(tagIDs, newTag.ID)
        } else {
            tagIDs = append(tagIDs, tag.ID)
        }
    }

    return tagIDs, nil
}

// UpdateTagUsageCount cập nhật số lần sử dụng của tag
func (s *TagService) UpdateTagUsageCount(tagIDs []uuid.UUID, increment bool) error {
    if len(tagIDs) == 0 {
        return nil
    }

    var change int64 = 1
    if !increment {
        change = -1
    }

    return config.DB.Model(&models.Tag{}).
        Where("id IN ?", tagIDs).
        UpdateColumn("usage_count", gorm.Expr("usage_count + ?", change)).Error
}

// SearchTags tìm kiếm tags theo tên
func (s *TagService) SearchTags(query string, limit int) ([]models.Tag, error) {
    var tags []models.Tag
    searchQuery := "%" + strings.ToLower(strings.TrimSpace(query)) + "%"
    
    if err := config.DB.Where("LOWER(name) LIKE ?", searchQuery).
        Order("usage_count DESC, name ASC").
        Limit(limit).
        Find(&tags).Error; err != nil {
        return nil, err
    }

    return tags, nil
} 