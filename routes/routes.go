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
    }

    return r
} 