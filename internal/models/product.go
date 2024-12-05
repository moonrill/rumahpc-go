package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
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
	Weight        float64        `gorm:"type:float;not null" validate:"required,gte=0" json:"weight"`
	Status        ProductStatus  `gorm:"type:varchar(255);default:'active';not null" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	MerchantID    string         `gorm:"type:uuid;not null" validate:"required" json:"merchant_id"`
	Merchant      *User          `json:"merchant" gorm:"foreignKey:MerchantID"`
	BrandID       *string        `gorm:"type:uuid;" json:"brand_id"`
	Brand         *Brand         `json:"brand" gorm:"foreignKey:BrandID"`
	CategoryID    string         `gorm:"type:uuid;not null" validate:"required" json:"category_id"`
	Category      *Category      `json:"category" gorm:"foreignKey:CategoryID"`
	SubCategory   *SubCategory   `json:"sub_category" gorm:"foreignKey:SubCategoryID"`
	SubCategoryID *string        `gorm:"type:uuid;" json:"sub_category_id"`
	Images        *[]string      `json:"images" gorm:"-"`
}

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
)

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	slug := utils.Slugify(p.Name)

	p.ID = uuid.New().String()
	p.Slug = fmt.Sprintf("%s-%s", slug, p.ID)
	return
}

func (p *Product) BeforeUpdate(tx *gorm.DB) (err error) {
	slug := utils.Slugify(p.Name)

	p.Slug = fmt.Sprintf("%s-%s", slug, p.ID)
	return
}
