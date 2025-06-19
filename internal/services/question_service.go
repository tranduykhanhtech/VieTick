package services

import (
    "errors"
    "time"

    "github.com/google/uuid"
    "vietick/config"
    "vietick/internal/models"
)

type QuestionService struct{}

type CreateQuestionRequest struct {
    Title   string   `json:"title" binding:"required,min=3"`
    Content string   `json:"content" binding:"required,min=3"`
    Tags    []string `json:"tags"` // Array of tag names
}

type UpdateQuestionRequest struct {
    Title   string   `json:"title" binding:"required,min=10"`
    Content string   `json:"content" binding:"required,min=20"`
    Tags    []string `json:"tags"` // Array of tag names
}

func NewQuestionService() *QuestionService {
    return &QuestionService{}
}

func (s *QuestionService) CreateQuestion(userID uuid.UUID, req CreateQuestionRequest) (*models.Question, error) {
    now := time.Now()
    question := models.Question{
        Title:     req.Title,
        Content:   req.Content,
        UserID:    userID,
        CreatedAt: now,
        UpdatedAt: now,
    }

    // Start transaction
    tx := config.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Create question
    if err := tx.Create(&question).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    // Handle tags if provided
    if len(req.Tags) > 0 {
        tagService := NewTagService()
        tagIDs, err := tagService.GetOrCreateTags(req.Tags)
        if err != nil {
            tx.Rollback()
            return nil, err
        }

        // Associate tags with question
        if len(tagIDs) > 0 {
            var tags []models.Tag
            if err := tx.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
                tx.Rollback()
                return nil, err
            }

            if err := tx.Model(&question).Association("Tags").Append(tags); err != nil {
                tx.Rollback()
                return nil, err
            }

            // Update usage count for tags
            if err := tagService.UpdateTagUsageCount(tagIDs, true); err != nil {
                tx.Rollback()
                return nil, err
            }
        }
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        return nil, err
    }

    // Load question with tags for response
    if err := config.DB.Preload("Tags").Preload("User").First(&question, question.ID).Error; err != nil {
        return nil, err
    }

    // Gửi notification đến followers về câu hỏi mới
    notificationService := NewNotificationService()
    notificationService.SendNotificationToFollowers(
        userID,
        models.NotificationTypeQuestion,
        "Câu hỏi mới từ người bạn follow",
        question.User.Username+" vừa đăng câu hỏi: "+question.Title,
        map[string]interface{}{
            "question_id": question.ID,
            "author_id": question.UserID,
            "author_name": question.User.Username,
        },
    )

    return &question, nil
}

func (s *QuestionService) GetQuestions(page, limit int) ([]models.Question, int64, error) {
    var questions []models.Question
    var total int64

    // Get total count
    if err := config.DB.Model(&models.Question{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Get questions with pagination and preload tags
    offset := (page - 1) * limit
    if err := config.DB.Preload("User").Preload("Tags").
        Order("created_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&questions).Error; err != nil {
        return nil, 0, err
    }

    return questions, total, nil
}

func (s *QuestionService) GetQuestionByID(questionID uuid.UUID) (*models.Question, error) {
    var question models.Question
    if err := config.DB.Preload("User").Preload("Tags").First(&question, "id = ?", questionID).Error; err != nil {
        return nil, errors.New("question not found")
    }
    return &question, nil
}

func (s *QuestionService) UpdateQuestion(questionID, userID uuid.UUID, req UpdateQuestionRequest) (*models.Question, error) {
    var question models.Question
    if err := config.DB.Preload("Tags").First(&question, "id = ?", questionID).Error; err != nil {
        return nil, errors.New("question not found")
    }

    // Check if user is the owner of the question
    if question.UserID != userID {
        return nil, errors.New("unauthorized to update this question")
    }

    // Start transaction
    tx := config.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Update basic fields
    question.Title = req.Title
    question.Content = req.Content
    question.UpdatedAt = time.Now()

    if err := tx.Save(&question).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    // Handle tags update
    tagService := NewTagService()
    
    // Get current tag IDs
    var currentTagIDs []uuid.UUID
    for _, tag := range question.Tags {
        currentTagIDs = append(currentTagIDs, tag.ID)
    }

    // Decrease usage count for current tags
    if len(currentTagIDs) > 0 {
        if err := tagService.UpdateTagUsageCount(currentTagIDs, false); err != nil {
            tx.Rollback()
            return nil, err
        }
    }

    // Remove all current tags
    if err := tx.Model(&question).Association("Tags").Clear(); err != nil {
        tx.Rollback()
        return nil, err
    }

    // Add new tags if provided
    if len(req.Tags) > 0 {
        newTagIDs, err := tagService.GetOrCreateTags(req.Tags)
        if err != nil {
            tx.Rollback()
            return nil, err
        }

        if len(newTagIDs) > 0 {
            var newTags []models.Tag
            if err := tx.Where("id IN ?", newTagIDs).Find(&newTags).Error; err != nil {
                tx.Rollback()
                return nil, err
            }

            if err := tx.Model(&question).Association("Tags").Append(newTags); err != nil {
                tx.Rollback()
                return nil, err
            }

            // Increase usage count for new tags
            if err := tagService.UpdateTagUsageCount(newTagIDs, true); err != nil {
                tx.Rollback()
                return nil, err
            }
        }
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        return nil, err
    }

    // Load updated question with tags
    if err := config.DB.Preload("Tags").Preload("User").First(&question, question.ID).Error; err != nil {
        return nil, err
    }

    return &question, nil
}

func (s *QuestionService) DeleteQuestion(questionID, userID uuid.UUID) error {
    var question models.Question
    if err := config.DB.Preload("Tags").First(&question, "id = ?", questionID).Error; err != nil {
        return errors.New("question not found")
    }

    // Check if user is the owner of the question
    if question.UserID != userID {
        return errors.New("unauthorized to delete this question")
    }

    // Start transaction
    tx := config.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Get tag IDs to decrease usage count
    var tagIDs []uuid.UUID
    for _, tag := range question.Tags {
        tagIDs = append(tagIDs, tag.ID)
    }

    // Delete all related answers first
    if err := tx.Where("question_id = ?", questionID).Delete(&models.Answer{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Delete all related votes
    if err := tx.Where("answer_id IN (SELECT id FROM answers WHERE question_id = ?)", questionID).Delete(&models.Vote{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Delete the question
    if err := tx.Delete(&question).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Decrease usage count for tags
    if len(tagIDs) > 0 {
        tagService := NewTagService()
        if err := tagService.UpdateTagUsageCount(tagIDs, false); err != nil {
            tx.Rollback()
            return err
        }
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        return err
    }

    return nil
}

// GetQuestionsByTag lấy câu hỏi theo tag
func (s *QuestionService) GetQuestionsByTag(tagName string, page, limit int) ([]models.Question, int64, error) {
    var questions []models.Question
    var total int64

    // Get total count
    if err := config.DB.Model(&models.Question{}).
        Joins("JOIN question_tags ON questions.id = question_tags.question_id").
        Joins("JOIN tags ON question_tags.tag_id = tags.id").
        Where("LOWER(tags.name) = LOWER(?)", tagName).
        Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Get questions with pagination
    offset := (page - 1) * limit
    if err := config.DB.Preload("User").Preload("Tags").
        Joins("JOIN question_tags ON questions.id = question_tags.question_id").
        Joins("JOIN tags ON question_tags.tag_id = tags.id").
        Where("LOWER(tags.name) = LOWER(?)", tagName).
        Order("questions.created_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&questions).Error; err != nil {
        return nil, 0, err
    }

    return questions, total, nil
}

// SearchQuestions tìm kiếm câu hỏi theo từ khóa
func (s *QuestionService) SearchQuestions(query string, page, limit int) ([]models.Question, int64, error) {
    var questions []models.Question
    var total int64

    searchQuery := "%" + query + "%"

    // Get total count
    if err := config.DB.Model(&models.Question{}).
        Where("LOWER(title) LIKE LOWER(?) OR LOWER(content) LIKE LOWER(?)", searchQuery, searchQuery).
        Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Get questions with pagination
    offset := (page - 1) * limit
    if err := config.DB.Preload("User").Preload("Tags").
        Where("LOWER(title) LIKE LOWER(?) OR LOWER(content) LIKE LOWER(?)", searchQuery, searchQuery).
        Order("created_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&questions).Error; err != nil {
        return nil, 0, err
    }

    return questions, total, nil
} 