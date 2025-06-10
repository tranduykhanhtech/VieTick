package utils

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

type Claims struct {
    UserID uuid.UUID `json:"user_id"`
    jwt.RegisteredClaims
}

func GenerateToken(userID uuid.UUID) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", fmt.Errorf("JWT_SECRET environment variable is not set")
    }

    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string) (*Claims, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return nil, fmt.Errorf("JWT_SECRET environment variable is not set")
    }

    log.Printf("Parsing token with secret length: %d", len(secret))
    
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })

    if err != nil {
        log.Printf("Token parsing error: %v", err)
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        log.Printf("Token is valid, user_id: %s", claims.UserID)
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token claims")
} 