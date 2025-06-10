package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Answer struct {
    ID         uuid.UUID  `gorm:"type:char(36);primaryKey;collate:utf8mb4_general_ci"`
    Content    string     `gorm:"type:text;not null;collate:utf8mb4_general_ci"`
    QuestionID uuid.UUID  `gorm:"type:char(36);not null;collate:utf8mb4_general_ci"`
    UserID     uuid.UUID  `gorm:"type:char(36);not null;collate:utf8mb4_general_ci"`
    IsVerified bool       `gorm:"default:false"`
    VerifiedBy *uuid.UUID `gorm:"type:char(36);collate:utf8mb4_general_ci"`
    CreatedAt  time.Time  `gorm:"not null"`
    UpdatedAt  time.Time  `gorm:"not null"`
    Reported   bool       `gorm:"default:false"`

    Question Question `gorm:"foreignKey:QuestionID;references:ID;constraint:OnDelete:CASCADE"`
    User     User     `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
    Verifier *User    `gorm:"foreignKey:VerifiedBy;references:ID"`
}

func (a *Answer) BeforeCreate(tx *gorm.DB) error {
    if a.ID == uuid.Nil {
        a.ID = uuid.New()
    }
    return nil
} 