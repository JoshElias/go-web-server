package internal

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewJwtToken(userId int, expirySeconds int) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirySeconds) * time.Second)),
		Subject:   strconv.Itoa(userId),
	}
	jwt := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	return jwt.SignedString(jwtSecret)
}
