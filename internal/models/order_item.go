package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderItem struct {
	ID        string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OrderID   string         `gorm:"type:uuid;not null" json:"order_id"`
	ProductID string         `gorm:"type:uuid;not null" json:"product_id"`
	Product   *Product       `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int            `gorm:"type:integer;not null" json:"quantity"`
	SubTotal  int            `gorm:"type:integer;not null" json:"sub_total"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
