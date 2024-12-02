package models

import (
	"time"

	"gorm.io/gorm"
)

type Brand struct {
	ID        string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" validate:"required,max=255" json:"name"`
	Slug      string         `gorm:"type:varchar(255);not null" json:"slug"`
	Image     string         `gorm:"type:text" json:"image"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Products  []Product      `json:"products" gorm:"foreignKey:BrandID"`
}
