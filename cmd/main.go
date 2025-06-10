package main

import (
    "vietick/config"
    "vietick/models"
    "vietick/routes"
)

func main() {
    config.InitDB()
    config.DB.AutoMigrate(&models.User{}, &models.Question{}, &models.Answer{}, &models.Vote{})
    r := routes.SetupRouter()
    r.Run(":8080")
} 