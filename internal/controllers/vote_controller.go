package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "vietick/internal/models"
    "vietick/internal/services"
)

type VoteController struct {
    voteService *services.VoteService
}

func NewVoteController(voteService *services.VoteService) *VoteController {
    return &VoteController{
        voteService: voteService,
    }
}

// VoteAnswer handles voting on an answer
func (c *VoteController) VoteAnswer(ctx *gin.Context) {
    // Get answer ID from URL parameter
    answerIDStr := ctx.Param("id")
    answerID, err := uuid.Parse(answerIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
        return
    }

    // Get user ID from context (set by auth middleware)
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Get vote type from URL parameter
    voteType := ctx.Param("type")
    if voteType != "up" && voteType != "down" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vote type"})
        return
    }

    // Create vote request
    req := services.CreateVoteRequest{
        Type: models.VoteType(voteType),
    }

    // Create vote
    vote, err := c.voteService.CreateVote(userID.(uuid.UUID), answerID, req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if vote == nil {
        ctx.JSON(http.StatusOK, gin.H{"message": "Vote removed"})
        return
    }

    ctx.JSON(http.StatusOK, vote)
}

// GetVotes handles getting votes for an answer
func (c *VoteController) GetVotes(ctx *gin.Context) {
    // Get answer ID from URL parameter
    answerIDStr := ctx.Param("id")
    answerID, err := uuid.Parse(answerIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
        return
    }

    // Get votes
    upVotes, downVotes, err := c.voteService.GetVotesByAnswer(answerID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "up_votes":   upVotes,
        "down_votes": downVotes,
    })
} 