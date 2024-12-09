package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ExternalID    string         `gorm:"type:varchar(255);not null" json:"external_id"`
	UserID        string         `gorm:"type:uuid;not null" json:"user_id"`
	User          User           `json:"user" gorm:"foreignKey:UserID"`
	Amount        int            `gorm:"type:integer;not null" json:"amount"`
	PaymentMethod string         `gorm:"type:varchar(255);" json:"payment_method"`
	PaymentDate   time.Time      `json:"payment_date"`
	Status        PaymentStatus  `gorm:"type:varchar(255);not null;default:'pending'" json:"status"`
	Orders        []Order        `json:"orders" gorm:"foreignKey:PaymentID"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "pending"
	PaymentPaid    PaymentStatus = "paid"
	PaymentFailed  PaymentStatus = "failed"
	PaymentRefund  PaymentStatus = "refund"
)
