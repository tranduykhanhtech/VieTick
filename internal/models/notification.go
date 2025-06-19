package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type NotificationType string

const (
    NotificationTypeFollow     NotificationType = "follow"
    NotificationTypeUnfollow   NotificationType = "unfollow"
    NotificationTypeAnswer     NotificationType = "answer"
    NotificationTypeVote       NotificationType = "vote"
    NotificationTypeVerify     NotificationType = "verify"
    NotificationTypeQuestion   NotificationType = "question"
    NotificationTypeTag        NotificationType = "tag"
)

type Notification struct {
    ID        uuid.UUID        `gorm:"type:char(36);primaryKey;collate:utf8mb4_general_ci"`
    UserID    uuid.UUID        `gorm:"type:char(36);not null;index;collate:utf8mb4_general_ci"` // Người nhận notification
    Type      NotificationType `gorm:"type:varchar(20);not null"`
    Title     string           `gorm:"type:varchar(255);not null"`
    Message   string           `gorm:"type:text;not null"`
    Data      string           `gorm:"type:json"` // JSON data cho additional info
    IsRead    bool             `gorm:"default:false"`
    CreatedAt time.Time        `gorm:"not null"`
    UpdatedAt time.Time        `gorm:"not null"`

    User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
    if n.ID == uuid.Nil {
        n.ID = uuid.New()
    }
    return nil
} 