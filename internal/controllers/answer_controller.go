package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "vietick/internal/services"
)

type AnswerController struct {
    answerService *services.AnswerService
}

func NewAnswerController() *AnswerController {
    return &AnswerController{
        answerService: services.NewAnswerService(),
    }
}

func (c *AnswerController) CreateAnswer(ctx *gin.Context) {
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

    questionIDStr := ctx.Param("id")
    questionID, err := uuid.Parse(questionIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid question ID"})
        return
    }

    var req services.CreateAnswerRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    answer, err := c.answerService.CreateAnswer(userIDUUID, questionID, req)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, answer)
}

func (c *AnswerController) GetAnswers(ctx *gin.Context) {
    questionIDStr := ctx.Param("id")
    questionID, err := uuid.Parse(questionIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid question ID"})
        return
    }

    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

    answers, total, err := c.answerService.GetAnswers(questionID, page, limit)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data":  answers,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

func (c *AnswerController) VerifyAnswer(ctx *gin.Context) {
    // Get user ID from context (verifier)
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Convert userID to uuid.UUID
    verifierID, ok := userID.(uuid.UUID)
    if !ok {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        return
    }

    answerIDStr := ctx.Param("id")
    answerID, err := uuid.Parse(answerIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
        return
    }

    if err := c.answerService.VerifyAnswer(answerID, verifierID); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Answer verified successfully"})
} 