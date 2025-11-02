package auth

import (
	"time"
	"github.com/dgrijalva/jwt-go"
)

var Secret = []byte("your-super-secret-jwt-key-2025-neftac-change-in-prod")

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func Generate(userID, role string) (string, error) {
	c := Claims{UserID: userID, Role: role, StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(24 * time.Hour).Unix()}}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString(Secret)
}
