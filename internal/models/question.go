package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Question struct {
    ID        uuid.UUID `gorm:"type:char(36);primaryKey;collate:utf8mb4_general_ci"`
    Title     string    `gorm:"type:varchar(255);not null;collate:utf8mb4_general_ci"`
    Content   string    `gorm:"type:text;not null;collate:utf8mb4_general_ci"`
    UserID    uuid.UUID `gorm:"type:char(36);not null;collate:utf8mb4_general_ci"`
    CreatedAt time.Time `gorm:"not null"`
    UpdatedAt time.Time `gorm:"not null"`

    User    User     `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
    Answers []Answer `gorm:"foreignKey:QuestionID;references:ID;constraint:OnDelete:CASCADE"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) error {
    if q.ID == uuid.Nil {
        q.ID = uuid.New()
    }
    return nil
} 