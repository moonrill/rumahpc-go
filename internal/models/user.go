package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" validate:"required,max=255" json:"name"`
	Email       string         `gorm:"uniqueIndex;type:varchar(255);not null" validate:"required,email,max=255" json:"email"`
	Password    string         `gorm:"type:varchar(255);not null" validate:"required,min=6,max=255" json:"password"`
	Avatar      *string        `gorm:"type:text" json:"avatar"`
	PhoneNumber string         `gorm:"type:varchar(13);not null" validate:"required,max=13" json:"phone_number"`
	Salt        string         `gorm:"type:uuid;not null" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
