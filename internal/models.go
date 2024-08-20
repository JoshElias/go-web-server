package internal

import "github.com/golang-jwt/jwt/v5"

type UserLoginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type UserLoginResponse struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

type UserDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserView struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type UserEntity struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

type ChirpDto struct {
	Body string `json:"body"`
}

type ChirpEntity struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

type RefreshToken struct {
	UserId    int              `json:"userId"`
	ExpiresAt *jwt.NumericDate `json:"expires_at"`
	Token     string           `json:"token"`
}
