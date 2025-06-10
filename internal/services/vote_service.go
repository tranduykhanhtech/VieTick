package services

import (
    "errors"
    "time"

    "github.com/google/uuid"
    "vietick/config"
    "vietick/internal/models"
)

type VoteService struct{}

type CreateVoteRequest struct {
    Type models.VoteType `json:"type" binding:"required,oneof=up down"`
}

const (
    VERIFICATION_THRESHOLD = 5 // Số upvote cần thiết để tự động xác minh
)

func NewVoteService() *VoteService {
    return &VoteService{}
}

func (s *VoteService) CreateVote(userID, answerID uuid.UUID, req CreateVoteRequest) (*models.Vote, error) {
    // Check if answer exists
    var answer models.Answer
    if err := config.DB.First(&answer, "id = ?", answerID).Error; err != nil {
        return nil, errors.New("answer not found")
    }

    // Check if user has already voted
    var existingVote models.Vote
    if err := config.DB.Where("user_id = ? AND answer_id = ?", userID, answerID).First(&existingVote).Error; err == nil {
        // If vote type is the same, remove the vote
        if existingVote.Type == req.Type {
            if err := config.DB.Delete(&existingVote).Error; err != nil {
                return nil, err
            }
            // Kiểm tra lại số upvote sau khi xóa vote
            if err := s.checkAndUpdateVerification(answerID); err != nil {
                return nil, err
            }
            return nil, nil
        }
        // If vote type is different, update the vote
        existingVote.Type = req.Type
        existingVote.UpdatedAt = time.Now()
        if err := config.DB.Save(&existingVote).Error; err != nil {
            return nil, err
        }
        // Kiểm tra lại số upvote sau khi thay đổi vote
        if err := s.checkAndUpdateVerification(answerID); err != nil {
            return nil, err
        }
        return &existingVote, nil
    }

    // Create new vote
    now := time.Now()
    vote := models.Vote{
        UserID:    userID,
        AnswerID:  answerID,
        Type:      req.Type,
        CreatedAt: now,
        UpdatedAt: now,
    }

    if err := config.DB.Create(&vote).Error; err != nil {
        return nil, err
    }

    // Kiểm tra số upvote sau khi tạo vote mới
    if err := s.checkAndUpdateVerification(answerID); err != nil {
        return nil, err
    }

    return &vote, nil
}

// checkAndUpdateVerification kiểm tra và cập nhật trạng thái xác minh của câu trả lời
func (s *VoteService) checkAndUpdateVerification(answerID uuid.UUID) error {
    var upVotes int64
    if err := config.DB.Model(&models.Vote{}).
        Where("answer_id = ? AND type = ?", answerID, models.UpVote).
        Count(&upVotes).Error; err != nil {
        return err
    }

    var answer models.Answer
    if err := config.DB.First(&answer, "id = ?", answerID).Error; err != nil {
        return err
    }

    // Nếu số upvote đạt ngưỡng và câu trả lời chưa được xác minh
    if upVotes >= VERIFICATION_THRESHOLD && !answer.IsVerified {
        answer.IsVerified = true
        // Lấy ID của người tạo câu trả lời làm người xác minh
        answer.VerifiedBy = &answer.UserID
        if err := config.DB.Save(&answer).Error; err != nil {
            return err
        }
    } else if upVotes < VERIFICATION_THRESHOLD && answer.IsVerified && answer.VerifiedBy != nil && *answer.VerifiedBy == answer.UserID {
        // Nếu số upvote giảm xuống dưới ngưỡng và câu trả lời đã được xác minh tự động
        answer.IsVerified = false
        answer.VerifiedBy = nil
        if err := config.DB.Save(&answer).Error; err != nil {
            return err
        }
    }

    return nil
}

func (s *VoteService) GetVotesByAnswer(answerID uuid.UUID) (int64, int64, error) {
    var upVotes, downVotes int64

    if err := config.DB.Model(&models.Vote{}).
        Where("answer_id = ? AND type = ?", answerID, models.UpVote).
        Count(&upVotes).Error; err != nil {
        return 0, 0, err
    }

    if err := config.DB.Model(&models.Vote{}).
        Where("answer_id = ? AND type = ?", answerID, models.DownVote).
        Count(&downVotes).Error; err != nil {
        return 0, 0, err
    }

    return upVotes, downVotes, nil
} 