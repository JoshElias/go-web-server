package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/JoshElias/go-web-server/internal"
	"github.com/JoshElias/go-web-server/internal/services"
	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("5")
	decoder := json.NewDecoder(r.Body)
	loginRequest := internal.UserLoginRequest{}
	err := decoder.Decode(&loginRequest)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	user, err := services.GetUserByEmail(loginRequest.Email)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	if err = bcrypt.CompareHashAndPassword(
		user.Password,
		[]byte(loginRequest.Password),
	); err != nil {
		fmt.Println("69")
		internal.RespondWithError(w, 401)
		return
	}
	accessToken, err := services.NewJwtToken(user.Id)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	refreshToken, err := services.NewRefreshToken(user.Id)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	internal.RespondWithJSON(
		w,
		200,
		internal.UserLoginResponse{
			Id:           user.Id,
			Email:        user.Email,
			IsChirpyRed:  user.IsChirpyRed,
			Token:        accessToken,
			RefreshToken: refreshToken.Token,
		},
	)
}

func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := getTokenFromAuthHeader(r)
	accessToken, err := services.RefreshAccessToken(refreshToken)
	if err != nil {
		if errors.Is(err, internal.RefreshTokenExpired) ||
			errors.Is(err, internal.RefreshTokenNotFound) {
			internal.RespondWithError(w, 401)
			return
		}
		internal.RespondWithError(w, 500)
		return
	}
	internal.RespondWithJSON(
		w,
		200,
		internal.RefreshTokenResponse{Token: accessToken},
	)
}

func HandleTokenRevoke(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromAuthHeader(r)
	_, err := services.RevokeRefreshToken(token)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	internal.RespondWithStatus(
		w,
		204,
	)
}

func getTokenFromAuthHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	return strings.TrimPrefix(authHeader, "Bearer ")
}
