package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Tag struct {
    ID          uuid.UUID `gorm:"type:char(36);primaryKey;collate:utf8mb4_general_ci"`
    Name        string    `gorm:"type:varchar(50);unique;not null;collate:utf8mb4_general_ci"`
    Description string    `gorm:"type:text;collate:utf8mb4_general_ci"`
    Color       string    `gorm:"type:varchar(7);default:'#007bff'"` // Hex color code
    UsageCount  int64     `gorm:"type:bigint;default:0"`             // Số lần sử dụng
    CreatedAt   time.Time `gorm:"not null"`
    UpdatedAt   time.Time `gorm:"not null"`

    Questions []Question `gorm:"many2many:question_tags;"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
    if t.ID == uuid.Nil {
        t.ID = uuid.New()
    }
    return nil
} 