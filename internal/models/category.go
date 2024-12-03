package models

import (
	"time"

	"github.com/moonrill/rumahpc-api/utils"
	"gorm.io/gorm"
)

type Category struct {
	ID            string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name          string         `gorm:"type:varchar(255);not null" validate:"required" json:"name"`
	Slug          string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"slug"`
	Icon          string         `gorm:"type:text" validate:"required" json:"icon"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	SubCategories *[]SubCategory `json:"sub_categories" gorm:"foreignKey:CategoryID"`
	Products      *[]Product     `json:"products" gorm:"foreignKey:CategoryID"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	c.Slug = utils.Slugify(c.Name)
	return
}

func (c *Category) BeforeUpdate(tx *gorm.DB) (err error) {
	c.Slug = utils.Slugify(c.Name)
	return
}
