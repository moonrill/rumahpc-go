package models

import (
	"time"

	"github.com/moonrill/rumahpc-api/utils"
	"gorm.io/gorm"
)

type Brand struct {
	ID        string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" validate:"required,max=255" json:"name"`
	Slug      string         `gorm:"type:varchar(255);" json:"slug"`
	Icon      string         `gorm:"type:text" json:"icon"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Products  []Product      `json:"products" gorm:"foreignKey:BrandID"`
}

func (b *Brand) BeforeCreate(tx *gorm.DB) (err error) {
	b.Slug = utils.Slugify(b.Name)
	return
}

func (b *Brand) BeforeUpdate(tx *gorm.DB) (err error) {
	b.Slug = utils.Slugify(b.Name)
	return
}
