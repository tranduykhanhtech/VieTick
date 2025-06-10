package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "vietick/internal/services"
)

type QuestionController struct {
    questionService *services.QuestionService
}

func NewQuestionController() *QuestionController {
    return &QuestionController{
        questionService: services.NewQuestionService(),
    }
}

func (c *QuestionController) CreateQuestion(ctx *gin.Context) {
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

    var req services.CreateQuestionRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    question, err := c.questionService.CreateQuestion(userIDUUID, req)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, question)
}

func (c *QuestionController) GetQuestions(ctx *gin.Context) {
    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

    questions, total, err := c.questionService.GetQuestions(page, limit)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data":  questions,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

func (c *QuestionController) GetQuestionByID(ctx *gin.Context) {
    questionIDStr := ctx.Param("id")
    questionID, err := uuid.Parse(questionIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid question ID"})
        return
    }

    question, err := c.questionService.GetQuestionByID(questionID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, question)
} 