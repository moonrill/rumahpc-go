package models

import (
	"time"

	"gorm.io/gorm"
)

type SubCategory struct {
	ID         string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CategoryID string         `gorm:"type:uuid;not null" json:"category_id"`
	Name       string         `gorm:"type:varchar(255);not null" json:"name"`
	Icon       string         `gorm:"type:text" json:"icon"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Category   Category       `json:"category"`
}
