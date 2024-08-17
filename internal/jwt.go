package internal

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewJwtToken(userId int) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(1) * time.Hour)),
		Subject:   strconv.Itoa(userId),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretString := os.Getenv("JWT_SECRET")
	if secretString == "" {
		return "", fmt.Errorf("JWT_SECRET env var not set")
	}
	return token.SignedString([]byte(secretString))
}
