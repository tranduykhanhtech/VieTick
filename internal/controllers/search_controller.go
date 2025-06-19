package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "vietick/internal/services"
)

type SearchController struct {
    questionService *services.QuestionService
    tagService      *services.TagService
}

func NewSearchController() *SearchController {
    return &SearchController{
        questionService: services.NewQuestionService(),
        tagService:      services.NewTagService(),
    }
}

// GetQuestionsByTag lấy câu hỏi theo tag
func (c *SearchController) GetQuestionsByTag(ctx *gin.Context) {
    tagName := ctx.Param("tag")
    if tagName == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "tag name is required"})
        return
    }

    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

    questions, total, err := c.questionService.GetQuestionsByTag(tagName, page, limit)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data":  questions,
        "total": total,
        "page":  page,
        "limit": limit,
        "tag":   tagName,
    })
}

// SearchQuestions tìm kiếm câu hỏi theo từ khóa
func (c *SearchController) SearchQuestions(ctx *gin.Context) {
    query := ctx.Query("q")
    if query == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "search query is required"})
        return
    }

    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

    questions, total, err := c.questionService.SearchQuestions(query, page, limit)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data":  questions,
        "total": total,
        "page":  page,
        "limit": limit,
        "query": query,
    })
}

// SearchTags tìm kiếm tags
func (c *SearchController) SearchTags(ctx *gin.Context) {
    query := ctx.Query("q")
    if query == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "search query is required"})
        return
    }

    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
    if limit > 50 {
        limit = 50
    }

    tags, err := c.tagService.SearchTags(query, limit)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"data": tags})
} 