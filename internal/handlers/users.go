package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/JoshElias/chirpy/internal"
	"github.com/JoshElias/chirpy/internal/services"
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
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}

	newUser, err := services.CreateUser(user.Email, hash)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	internal.RespondWithJSON(w, 201, newUser)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userLogin := internal.UserDto{}
	err := decoder.Decode(&userLogin)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}

	user, err := services.GetUserByEmail(userLogin.Email)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	err = bcrypt.CompareHashAndPassword(
		user.Password,
		[]byte(userLogin.Password),
	)
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
