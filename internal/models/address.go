package models

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID            string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID        string         `gorm:"type:uuid;not null" json:"user_id"`
	ContactName   string         `gorm:"type:varchar(255);not null" validate:"required,max=255" json:"contact_name"`
	ContactNumber string         `gorm:"type:varchar(13);not null" validate:"required,max=13" json:"contact_number"`
	Province      string         `gorm:"type:varchar(255);not null" validate:"required,max=255" json:"province"`
	City          string         `gorm:"type:varchar(255);not null" validate:"required,max=255" json:"city"`
	District      string         `gorm:"type:varchar(255);not null" validate:"required,max=255" json:"district"`
	Village       string         `gorm:"type:varchar(255);not null" validate:"required,max=255" json:"village"`
	ZipCode       string         `gorm:"type:varchar(10);not null" validate:"required,max=10" json:"zip_code"`
	Address       string         `gorm:"type:text;not null" validate:"required" json:"address"`
	Note          string         `gorm:"type:text" json:"note"`
	Default       bool           `json:"default" gorm:"default:false" validate:"boolean"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
