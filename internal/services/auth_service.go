package services

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/moonrill/rumahpc-api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Sub    string `json:"sub"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, string, error) {
	salt := uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)

	return string(hashedPassword), salt, err
}

func GenerateToken(user *models.User) (string, error) {
	jwtSecretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	var avatar string
	if user.Avatar != nil {
		avatar = string(*user.Avatar)
	} else {
		avatar = ""
	}

	claims := Claims{
		user.ID.String(),
		user.Name,
		avatar,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "rumahpc-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
