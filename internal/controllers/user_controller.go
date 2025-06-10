package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "vietick/internal/services"
)

type UserController struct {
    userService *services.UserService
}

func NewUserController() *UserController {
    return &UserController{
        userService: services.NewUserService(),
    }
}

func (c *UserController) Register(ctx *gin.Context) {
    var req services.RegisterRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response, err := c.userService.Register(req)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, response)
}

func (c *UserController) Login(ctx *gin.Context) {
    var req services.LoginRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response, err := c.userService.Login(req)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, response)
}

func (c *UserController) GetProfile(ctx *gin.Context) {
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

    user, err := c.userService.GetProfile(userIDUUID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, user)
} 