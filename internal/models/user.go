package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Email       string         `gorm:"uniqueIndex;type:varchar(255);not null" json:"email"`
	Password    string         `gorm:"type:varchar(255);not null" json:"-"`
	Avatar      *string        `gorm:"type:text" json:"avatar"`
	PhoneNumber string         `gorm:"type:varchar(13);not null" json:"phone_number"`
	Salt        string         `gorm:"type:uuid;not null" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	RoleID      string         `gorm:"type:uuid;not null" json:"role_id"`
	Role        Role           `json:"role" gorm:"foreignKey:RoleID"`
	Addresses   *[]Address     `json:"addresses" gorm:"foreignKey:UserID"`
	Products    *[]Product     `json:"products,omitempty" gorm:"foreignKey:MerchantID"`
	Orders      *[]Order       `json:"orders" gorm:"foreignKey:UserID"`
	Payments    *[]Payment     `json:"payments" gorm:"foreignKey:UserID"`
	Cart        *Cart          `json:"cart" gorm:"foreignKey:UserID"`
}
