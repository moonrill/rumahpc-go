package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    string         `gorm:"type:uuid;not null" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	CartItems []CartItem     `json:"cart_items" gorm:"foreignKey:CartID"`
}
