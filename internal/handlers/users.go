package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/JoshElias/chirpy/internal"
	"github.com/JoshElias/chirpy/internal/services"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HandleAddUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := internal.UserDto{}
	err := decoder.Decode(&user)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}

	newUser, err := services.CreateUser(user)
	if err != nil {
		if errors.Is(err, internal.UserAlreadyExists) {
			internal.RespondWithError(w, 400)
			return
		}
		internal.RespondWithError(w, 500)
		return
	}
	internal.RespondWithJSON(w, 201, newUser)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
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
	err = bcrypt.CompareHashAndPassword(
		user.Password,
		[]byte(loginRequest.Password),
	)
	if err != nil {
		internal.RespondWithError(w, 401)
		return
	}
	token, err := internal.NewJwtToken(user.Id)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	internal.RespondWithJSON(
		w,
		200,
		internal.UserLoginResponse{
			Id:    user.Id,
			Email: user.Email,
			Token: token,
		},
	)
}

func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		internal.RespondWithError(w, 401)
		return
	}
	userIdString, err := token.Claims.GetSubject()
	if err != nil {
		internal.RespondWithError(w, 401)
		return
	}
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}

	decoder := json.NewDecoder(r.Body)
	userPatch := internal.UserDto{}
	err = decoder.Decode(&userPatch)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}

	user, err := services.UpdateUserById(userId, userPatch)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}

	internal.RespondWithJSON(
		w,
		200,
		internal.UserView{
			Id:    user.Id,
			Email: user.Email,
		},
	)
}
