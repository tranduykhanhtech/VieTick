package routes

import (
    "github.com/gin-gonic/gin"
    "vietick/internal/controllers"
    "vietick/internal/middleware"
    "vietick/internal/services"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // Add CORS middleware
    r.Use(middleware.CORSMiddleware())

    // Initialize services
    voteService := services.NewVoteService()

    // Initialize controllers
    userController := controllers.NewUserController()
    questionController := controllers.NewQuestionController()
    answerController := controllers.NewAnswerController()
    voteController := controllers.NewVoteController(voteService)
    tagController := controllers.NewTagController()
    searchController := controllers.NewSearchController()
    followController := controllers.NewFollowController()

    // Public routes
    r.POST("/register", userController.Register)
    r.POST("/login", userController.Login)

    // Protected routes
    protected := r.Group("/")
    protected.Use(middleware.AuthMiddleware())
    {
        // User routes
        protected.GET("/users/me", userController.GetProfile)

        // Question routes
        protected.POST("/questions", questionController.CreateQuestion)
        protected.GET("/questions", questionController.GetQuestions)
        protected.GET("/questions/:id", questionController.GetQuestionByID)
        protected.PUT("/questions/:id", questionController.UpdateQuestion)
        protected.DELETE("/questions/:id", questionController.DeleteQuestion)

        // Answer routes
        protected.POST("/questions/:id/answers", answerController.CreateAnswer)
        protected.GET("/questions/:id/answers", answerController.GetAnswers)

        // Answer group for specific answer operations
        answerGroup := protected.Group("/answers")
        {
            answerIDGroup := answerGroup.Group("/:id")
            {
                answerIDGroup.POST("/verify", answerController.VerifyAnswer)    // /answers/:id/verify
                answerIDGroup.POST("/vote/:type", voteController.VoteAnswer)    // /answers/:id/vote/:type
                answerIDGroup.GET("/votes", voteController.GetVotes)            // /answers/:id/votes
            }
        }

        // Tag management routes
        tagGroup := protected.Group("/tags")
        {
            tagGroup.POST("", tagController.CreateTag)                    // POST /tags
            tagGroup.GET("", tagController.GetTags)                       // GET /tags
            tagGroup.GET("/:id", tagController.GetTagByID)                // GET /tags/:id
            tagGroup.PUT("/:id", tagController.UpdateTag)                 // PUT /tags/:id
            tagGroup.DELETE("/:id", tagController.DeleteTag)              // DELETE /tags/:id
        }

        // Search routes
        searchGroup := protected.Group("/search")
        {
            searchGroup.GET("/questions", searchController.SearchQuestions)     // GET /search/questions?q=query
            searchGroup.GET("/tags", searchController.SearchTags)               // GET /search/tags?q=query
            searchGroup.GET("/questions/tag/:tag", searchController.GetQuestionsByTag) // GET /search/questions/tag/:tag
        }

        // Follow routes
        followGroup := protected.Group("/follows")
        {
            followGroup.POST("/:id", followController.FollowUser)              // POST /follows/:id (follow user)
            followGroup.DELETE("/:id", followController.UnfollowUser)          // DELETE /follows/:id (unfollow user)
            followGroup.GET("/:id/check", followController.IsFollowing)        // GET /follows/:id/check (check if following)
            followGroup.GET("/:id/followers", followController.GetFollowers)   // GET /follows/:id/followers (get user's followers)
            followGroup.GET("/:id/following", followController.GetFollowing)   // GET /follows/:id/following (get who user is following)
            followGroup.GET("/:id/stats", followController.GetUserFollowStats) // GET /follows/:id/stats (get user's follow stats)
        }

        // My follow stats
        protected.GET("/me/follows/stats", followController.GetMyFollowStats)  // GET /me/follows/stats (get my follow stats)
    }

    return r
} 