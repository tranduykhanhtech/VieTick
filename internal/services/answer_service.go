package services

import (
	"errors"
	"log"
	"time"

	"vietick/config"
	"vietick/internal/models"

	// "vietick/internal/services"

	"github.com/google/uuid"
)

type AnswerService struct{}

type CreateAnswerRequest struct {
	Content string `json:"content" binding:"required,min=10"`
}

func NewAnswerService() *AnswerService {
	return &AnswerService{}
}

func (s *AnswerService) CreateAnswer(userID, questionID uuid.UUID, req CreateAnswerRequest) (*models.Answer, error) {
	log.Printf("Creating answer for user ID: %s, question ID: %s", userID, questionID)

	// Check if question exists
	var question models.Question
	if err := config.DB.First(&question, "id = ?", questionID).Error; err != nil {
		log.Printf("Question not found: %v", err)
		return nil, errors.New("question not found")
	}
	log.Printf("Question found: %s", question.ID)

	// Check if user exists
	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		log.Printf("User not found: %v", err)
		// Log all users in database for debugging
		var users []models.User
		if err := config.DB.Find(&users).Error; err != nil {
			log.Printf("Error fetching users: %v", err)
		} else {
			log.Printf("Available users in database:")
			for _, u := range users {
				log.Printf("User ID: %s, Email: %s, Username: %s", u.ID, u.Email, u.Username)
			}
		}
		return nil, errors.New("user not found")
	}
	log.Printf("User found: %s", user.ID)

	now := time.Now()
	answer := models.Answer{
		Content:    req.Content,
		UserID:     userID,
		QuestionID: questionID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := config.DB.Create(&answer).Error; err != nil {
		log.Printf("Error creating answer: %v", err)
		return nil, err
	}

	log.Printf("Answer created successfully with ID: %s", answer.ID)

	return &answer, nil
}

func (s *AnswerService) GetAnswers(questionID uuid.UUID, page, limit int) ([]models.Answer, int64, error) {
	var answers []models.Answer
	var total int64

	// Get total count
	if err := config.DB.Model(&models.Answer{}).
		Where("question_id = ?", questionID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get answers with pagination
	offset := (page - 1) * limit
	if err := config.DB.Preload("User").
		Where("question_id = ?", questionID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&answers).Error; err != nil {
		return nil, 0, err
	}

	return answers, total, nil
}

func (s *AnswerService) VerifyAnswer(answerID, verifierID uuid.UUID) error {
	var answer models.Answer
	if err := config.DB.First(&answer, "id = ?", answerID).Error; err != nil {
		return errors.New("answer not found")
	}

	// Nếu câu trả lời đã được xác minh, bỏ xác minh
	if answer.IsVerified {
		answer.IsVerified = false
		answer.VerifiedBy = nil
	} else {
		// Nếu chưa được xác minh, thêm xác minh
		answer.IsVerified = true
		answer.VerifiedBy = &verifierID
	}

	if err := config.DB.Save(&answer).Error; err != nil {
		return err
	}

	return nil
}
