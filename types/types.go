package types

import "github.com/golang-jwt/jwt/v5"

type SignUpRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Email       string `json:"email" validate:"required,email,max=255"`
	Password    string `json:"password" validate:"required,min=8,max=255"`
	PhoneNumber string `json:"phone_number" validate:"required,max=13"`
	Role        string `json:"role" validate:"required,max=255"`
}

type JwtClaims struct {
	Sub    string `json:"sub"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	jwt.RegisteredClaims
}
