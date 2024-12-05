package types

import "github.com/golang-jwt/jwt/v5"

type SignUpRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Email       string `json:"email" validate:"required,email,max=255"`
	Password    string `json:"password" validate:"required,min=8,max=255"`
	PhoneNumber string `json:"phone_number" validate:"required,min=10,max=13"`
	Role        string `json:"role" validate:"required,max=255"`
}

type JwtClaims struct {
	Sub    string `json:"sub"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type CreateProductRequest struct {
	Name          string   `json:"name" validate:"required,max=255"`
	Description   string   `json:"description" validate:"required"`
	Stock         int      `json:"stock" validate:"required,gte=0"`
	Price         int      `json:"price" validate:"required,gte=0"`
	Weight        float64  `json:"weight" validate:"required,gte=0"`
	CategoryID    string   `json:"category_id" validate:"required"`
	SubCategoryID *string  `json:"sub_category_id"`
	BrandID       *string  `json:"brand_id"`
	Images        []string `json:"images" `
}

type UpdateProductRequest struct {
	Name          string  `json:"name" validate:"required,max=255"`
	Description   string  `json:"description" validate:"required"`
	Stock         int     `json:"stock" validate:"required,gte=0"`
	Price         int     `json:"price" validate:"required,gte=0"`
	Weight        float64 `json:"weight" validate:"required,gte=0"`
	CategoryID    string  `json:"category_id" validate:"required"`
	SubCategoryID *string `json:"sub_category_id"`
	BrandID       *string `json:"brand_id"`
}

type AddToCartRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gte=1"`
}

type UpdateCartRequest struct {
	Quantity int `json:"quantity" validate:"required,gte=1"`
}

type RemoveFromCartRequest struct {
	CartItemsID []string `json:"cart_items_id" validate:"required,gte=1"`
}
