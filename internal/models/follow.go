package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Follow struct {
    ID          uuid.UUID `gorm:"type:char(36);primaryKey;collate:utf8mb4_general_ci"`
    FollowerID  uuid.UUID `gorm:"type:char(36);not null;index;collate:utf8mb4_general_ci"` // Người follow
    FollowingID uuid.UUID `gorm:"type:char(36);not null;index;collate:utf8mb4_general_ci"` // Người được follow
    CreatedAt   time.Time `gorm:"not null"`

    Follower  User `gorm:"foreignKey:FollowerID;references:ID;constraint:OnDelete:CASCADE"`
    Following User `gorm:"foreignKey:FollowingID;references:ID;constraint:OnDelete:CASCADE"`
}

func (f *Follow) BeforeCreate(tx *gorm.DB) error {
    if f.ID == uuid.Nil {
        f.ID = uuid.New()
    }
    return nil
} 