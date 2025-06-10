package services

import (
    "errors"
    "log"
    "time"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "vietick/config"
    "vietick/internal/models"
    "vietick/pkg/utils"
)

type UserService struct{}

type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Username string `json:"username" binding:"required,min=3,max=20"`
    Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Token string      `json:"token"`
    User  models.User `json:"user"`
}

func NewUserService() *UserService {
    return &UserService{}
}

func (s *UserService) Register(req RegisterRequest) (*LoginResponse, error) {
    // Check if email exists
    var existingUser models.User
    if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
        return nil, errors.New("email already exists")
    }

    // Check if username exists
    if err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
        return nil, errors.New("username already exists")
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // Create user
    now := time.Now()
    user := models.User{
        Email:     req.Email,
        Username:  req.Username,
        Password:  string(hashedPassword),
        Point:     0,
        CreatedAt: now,
        UpdatedAt: now,
    }

    if err := config.DB.Create(&user).Error; err != nil {
        return nil, err
    }
    log.Printf("User registered with ID: %s", user.ID)

    // Generate JWT token
    token, err := utils.GenerateToken(user.ID)
    if err != nil {
        return nil, err
    }

    return &LoginResponse{
        Token: token,
        User:  user,
    }, nil
}

func (s *UserService) Login(req LoginRequest) (*LoginResponse, error) {
    var user models.User
    if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        return nil, errors.New("invalid email or password")
    }

    // Check password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, errors.New("invalid email or password")
    }

    // Generate JWT token
    token, err := utils.GenerateToken(user.ID)
    if err != nil {
        return nil, err
    }

    return &LoginResponse{
        Token: token,
        User:  user,
    }, nil
}

func (s *UserService) GetProfile(userID uuid.UUID) (*models.User, error) {
    log.Printf("Attempting to retrieve user with ID: %s", userID)
    var user models.User
    if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
        return nil, errors.New("user not found")
    }
    return &user, nil
}

func (s *UserService) AddPoint(userID uuid.UUID, points int) error {
    return config.DB.Model(&models.User{}).Where("id = ?", userID).
        UpdateColumn("point", gorm.Expr("point + ?", points)).Error
} 