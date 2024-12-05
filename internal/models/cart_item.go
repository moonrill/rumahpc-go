package models

import (
	"time"

	"gorm.io/gorm"
)

type CartItem struct {
	ID        string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CartID    string         `gorm:"type:uuid;not null" json:"cart_id"`
	ProductID string         `gorm:"type:uuid;not null" json:"product_id"`
	Product   *Product       `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int            `gorm:"type:integer;not null" json:"quantity"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
