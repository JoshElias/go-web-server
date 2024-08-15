package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/JoshElias/chirpy/internal"
)

func HandleAddUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := internal.UserDto{}
	err := decoder.Decode(&user)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}

	conn, err := internal.GetTestDbConnection()
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
