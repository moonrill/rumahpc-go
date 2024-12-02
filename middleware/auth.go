package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
)

func CookiesAuthMiddleware(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	jwtSecretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "status_code": http.StatusUnauthorized})
		return
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}

		return jwtSecretKey, nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "status_code": http.StatusUnauthorized})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user models.User
		err := config.DB.Preload("Role").First(&user, "id = ?", claims["sub"]).Error

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "status_code": http.StatusUnauthorized})
			return
		}

		c.Set("user", user)
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "status_code": http.StatusUnauthorized})
	}
}

func HeaderAuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	jwtSecretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "status_code": http.StatusUnauthorized, "message": "Authorization header missing"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "status_code": http.StatusUnauthorized, "message": "Invalid authorization header"})
		return
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}

		return jwtSecretKey, nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "status_code": http.StatusUnauthorized})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user models.User
		err := config.DB.Preload("Role").First(&user, "id = ?", claims["sub"]).Error

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "status_code": http.StatusUnauthorized})
			return
		}

		c.Set("user", user)
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "status_code": http.StatusUnauthorized})
	}
}
