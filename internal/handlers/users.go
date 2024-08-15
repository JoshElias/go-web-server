package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/JoshElias/chirpy/internal"
)

type UserDto struct {
	Email string `json:"email"`
}

type UserEntity struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func HandleAddUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := UserDto{}
	err := decoder.Decode(&user)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}

	conn, err := getDbConnection()
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	newChirp, err := conn.CreateUser(user.Email)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	internal.RespondWithJSON(w, 201, newChirp)
}
