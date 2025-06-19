package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "vietick/internal/services"
)

type TagController struct {
    tagService *services.TagService
}

func NewTagController() *TagController {
    return &TagController{
        tagService: services.NewTagService(),
    }
}

// CreateTag tạo tag mới
func (c *TagController) CreateTag(ctx *gin.Context) {
    var req services.CreateTagRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    tag, err := c.tagService.CreateTag(req)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, tag)
}

// GetTags lấy danh sách tags
func (c *TagController) GetTags(ctx *gin.Context) {
    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

    tags, total, err := c.tagService.GetTags(page, limit)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data":  tags,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

// GetTagByID lấy tag theo ID
func (c *TagController) GetTagByID(ctx *gin.Context) {
    tagIDStr := ctx.Param("id")
    tagID, err := uuid.Parse(tagIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag ID"})
        return
    }

    tag, err := c.tagService.GetTagByID(tagID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, tag)
}

// UpdateTag cập nhật tag
func (c *TagController) UpdateTag(ctx *gin.Context) {
    tagIDStr := ctx.Param("id")
    tagID, err := uuid.Parse(tagIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag ID"})
        return
    }

    var req services.UpdateTagRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    tag, err := c.tagService.UpdateTag(tagID, req)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, tag)
}

// DeleteTag xóa tag
func (c *TagController) DeleteTag(ctx *gin.Context) {
    tagIDStr := ctx.Param("id")
    tagID, err := uuid.Parse(tagIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag ID"})
        return
    }

    if err := c.tagService.DeleteTag(tagID); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "tag deleted successfully"})
}

// SearchTags tìm kiếm tags
func (c *TagController) SearchTags(ctx *gin.Context) {
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