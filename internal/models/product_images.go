package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductImages struct {
	ID        string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProductID string         `gorm:"type:uuid;not null" json:"product_id"`
	Product   *Product       `json:"product" gorm:"foreignKey:ProductID"`
	Image     string         `gorm:"type:text;not null" json:"image"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
