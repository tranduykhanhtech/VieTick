package main

import (
    "log"
    "os"

    "vietick/config"
    "vietick/internal/models"
    "vietick/routes"
    "github.com/gin-gonic/gin"
)

func main() {
    gin.SetMode(gin.DebugMode)
    // Initialize database
    if err := config.InitDB(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer config.CloseDB()

    // Check if we need to reset database
    if os.Getenv("RESET_DB") == "true" {
        log.Println("Resetting database...")
        // Drop existing tables
        if err := config.DB.Migrator().DropTable(&models.Vote{}, &models.Answer{}, &models.Question{}, &models.User{}); err != nil {
            log.Fatalf("Failed to drop tables: %v", err)
        }
    }

    // Auto migrate database schema
    if err := config.DB.AutoMigrate(
        &models.User{},
        &models.Question{},
        &models.Answer{},
        &models.Vote{},
    ); err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }

    // Setup router
    r := routes.SetupRouter()

    // Start server
    log.Println("Server starting on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
} 