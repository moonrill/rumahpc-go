package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Email       string         `gorm:"uniqueIndex;type:varchar(255);not null" json:"email"`
	Password    string         `gorm:"type:varchar(255);not null" json:"-"`
	Avatar      *string        `gorm:"type:text" json:"avatar"`
	PhoneNumber string         `gorm:"type:varchar(13);not null" json:"phone_number"`
	Salt        string         `gorm:"type:uuid;not null" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
