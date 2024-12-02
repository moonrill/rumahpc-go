package models

import (
	"time"

	"github.com/moonrill/rumahpc-api/utils"
	"gorm.io/gorm"
)

type Product struct {
	ID            string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name          string         `gorm:"type:varchar(255);not null" validate:"required,max=255" json:"name"`
	Slug          string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"slug"`
	Description   string         `gorm:"type:text" validate:"required" json:"description"`
	Price         int            `gorm:"type:integer;not null" validate:"required,gte=0" json:"price"`
	Stock         int            `gorm:"type:integer;not null" validate:"required,gte=0" json:"stock"`
	Weight        int            `gorm:"type:integer;not null" validate:"required,gte=0" json:"weight"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	BrandID       *string        `gorm:"type:uuid;" json:"brand_id"`
	Brand         *Brand         `json:"brand" gorm:"foreignKey:BrandID"`
	CategoryID    string         `gorm:"type:uuid;not null" validate:"required" json:"category_id"`
	Category      Category       `json:"category" gorm:"foreignKey:CategoryID"`
	SubCategory   *SubCategory   `json:"sub_category" gorm:"foreignKey:SubCategoryID"`
	SubCategoryID *string        `gorm:"type:uuid;" json:"sub_category_id"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.Slug = utils.Slugify(p.Name)
	return
}

func (p *Product) BeforeUpdate(tx *gorm.DB) (err error) {
	p.Slug = utils.Slugify(p.Name)
	return
}
