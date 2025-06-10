package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type VoteType string

const (
    UpVote   VoteType = "up"
    DownVote VoteType = "down"
)

type Vote struct {
    ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
    UserID    uuid.UUID `gorm:"type:char(36);index;not null"`
    AnswerID  uuid.UUID `gorm:"type:char(36);index;not null"`
    Type      VoteType  `gorm:"type:varchar(10);not null"`
    CreatedAt time.Time `gorm:"not null"`
    UpdatedAt time.Time `gorm:"not null"`

    User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
    Answer Answer `gorm:"foreignKey:AnswerID;references:ID;constraint:OnDelete:CASCADE"`
}

func (v *Vote) BeforeCreate(tx *gorm.DB) error {
    if v.ID == uuid.Nil {
        v.ID = uuid.New()
    }
    return nil
} 