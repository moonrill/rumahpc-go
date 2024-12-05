package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID        string         `gorm:"type:uuid;not null" json:"user_id"`
	User          User           `json:"user" gorm:"foreignKey:UserID"`
	OrderID       string         `gorm:"type:uuid;not null" json:"order_id"`
	Order         Order          `json:"order" gorm:"foreignKey:OrderID"`
	Amount        int            `gorm:"type:integer;not null" json:"amount"`
	PaymentMethod string         `gorm:"type:varchar(255);" json:"payment_method"`
	PaymentDate   time.Time      `json:"payment_date"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
