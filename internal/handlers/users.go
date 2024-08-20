package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JoshElias/go-web-server/internal"
	"github.com/JoshElias/go-web-server/internal/services"
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

func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		internal.RespondWithError(w, 401)
		return
	}
	decoder := json.NewDecoder(r.Body)
	userPatch := internal.UserDto{}
	err := decoder.Decode(&userPatch)
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
