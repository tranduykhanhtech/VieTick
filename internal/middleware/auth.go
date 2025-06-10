package middleware

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
    "vietick/pkg/utils"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Printf("JWT_SECRET from env: %s", os.Getenv("JWT_SECRET"))
        
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            log.Printf("Missing Authorization header")
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        log.Printf("Token received: %s", tokenString)
        
        claims, err := utils.ParseToken(tokenString)
        if err != nil {
            log.Printf("Token parsing error: %v", err)
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
            return
        }

        log.Printf("Token parsed successfully, user_id: %s", claims.UserID)
        c.Set("user_id", claims.UserID)
        c.Next()
    }
} 