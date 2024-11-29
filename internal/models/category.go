package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID            string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name          string         `gorm:"type:varchar(255);not null" validate:"required" json:"name"`
	Icon          string         `gorm:"type:text" validate:"required" json:"icon"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	SubCategories []SubCategory  `json:"sub_categories" gorm:"foreignKey:CategoryID"`
}
