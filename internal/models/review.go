package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID          string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      string         `gorm:"type:uuid;not null" json:"user_id"`
	MerchantID  string         `gorm:"type:uuid;not null" json:"merchant_id"`
	User        *User          `json:"user" gorm:"foreignKey:UserID"`
	ProductID   string         `gorm:"type:uuid;not null" json:"product_id"`
	Product     *Product       `json:"product" gorm:"foreignKey:ProductID"`
	OrderItemID string         `gorm:"type:uuid;not null;unique" json:"order_item_id"`
	OrderItem   *OrderItem     `json:"order_item" gorm:"foreignKey:OrderItemID"`
	Rating      int            `gorm:"type:integer;not null" json:"rating"`
	Comment     *string        `gorm:"type:text" json:"comment"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
