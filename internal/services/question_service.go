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
    Title   string `json:"title" binding:"required,min=10"`
    Content string `json:"content" binding:"required,min=20"`
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

    if err := config.DB.Create(&question).Error; err != nil {
        return nil, err
    }

    return &question, nil
}

func (s *QuestionService) GetQuestions(page, limit int) ([]models.Question, int64, error) {
    var questions []models.Question
    var total int64

    // Get total count
    if err := config.DB.Model(&models.Question{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Get questions with pagination
    offset := (page - 1) * limit
    if err := config.DB.Preload("User").
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
    if err := config.DB.Preload("User").First(&question, "id = ?", questionID).Error; err != nil {
        return nil, errors.New("question not found")
    }
    return &question, nil
} 