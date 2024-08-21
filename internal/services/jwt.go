package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/JoshElias/go-web-server/internal"
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

func NewRefreshToken(userId int) (internal.RefreshToken, error) {
	randomData := make([]byte, 32)
	_, err := rand.Read(randomData)
	if err != nil {
		return internal.RefreshToken{}, nil
	}
	tokenString := hex.EncodeToString(randomData)
	refreshToken := internal.RefreshToken{
		UserId:    userId,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(60*24) * time.Hour)),
		Token:     tokenString,
	}
	conn, err := internal.GetTestDbConnection()
	if err != nil {
		return internal.RefreshToken{}, err
	}
	db, err := conn.LoadDb()
	if err != nil {
		return internal.RefreshToken{}, err
	}
	_, exists := db.RefreshTokens[tokenString]
	if exists {
		return internal.RefreshToken{}, fmt.Errorf("we got a collision in access tokens")
	}
	db.RefreshTokens[tokenString] = refreshToken
	err = conn.WriteDb(db)
	if err != nil {
		return internal.RefreshToken{}, err
	}
	return refreshToken, nil
}

func RefreshAccessToken(refreshToken string) (string, error) {
	conn, err := internal.GetTestDbConnection()
	if err != nil {
		return "", err
	}
	db, err := conn.LoadDb()
	if err != nil {
		return "", err
	}
	tokenEntity, exists := db.RefreshTokens[refreshToken]
	if !exists {
		return "", internal.RefreshTokenNotFound
	}
	if tokenEntity.ExpiresAt.Time.Before(time.Now()) {
		return "", internal.RefreshTokenExpired
	}

	newToken, err := NewJwtToken(tokenEntity.UserId)
	if err != nil {
		return "", err
	}
	return newToken, nil
}

func RevokeRefreshToken(token string) (bool, error) {
	conn, err := internal.GetTestDbConnection()
	if err != nil {
		return false, err
	}
	db, err := conn.LoadDb()
	if err != nil {
		return false, err
	}
	_, exists := db.RefreshTokens[token]
	if !exists {
		return false, nil
	}
	delete(db.RefreshTokens, token)
	err = conn.WriteDb(db)
	return true, nil
}
