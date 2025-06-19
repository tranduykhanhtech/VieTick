package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID        uuid.UUID `gorm:"type:char(36);primaryKey;collate:utf8mb4_general_ci"`
    Email     string    `gorm:"type:varchar(255);unique;not null;collate:utf8mb4_general_ci"`
    Username  string    `gorm:"type:varchar(50);unique;not null;collate:utf8mb4_general_ci"`
    Password  string    `gorm:"type:varchar(255);not null;collate:utf8mb4_general_ci"`
    Point     int64     `gorm:"type:bigint;default:0"`
    CreatedAt time.Time `gorm:"not null"`
    UpdatedAt time.Time `gorm:"not null"`

    Questions []Question `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
    Answers   []Answer   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
    Votes     []Vote     `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
    
    // Follow relationships
    Followers  []Follow `gorm:"foreignKey:FollowingID;references:ID;constraint:OnDelete:CASCADE"` // Những người follow mình
    Following  []Follow `gorm:"foreignKey:FollowerID;references:ID;constraint:OnDelete:CASCADE"`  // Những người mình follow
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return nil
} 