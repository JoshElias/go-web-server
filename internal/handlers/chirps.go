package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/JoshElias/go-web-server/internal"
	"github.com/JoshElias/go-web-server/internal/services"
)

var badWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

var badWordTree *internal.Trie

func init() {
	badWordTree = internal.NewTrie()
	for _, w := range badWords {
		badWordTree.Add(w)
	}
}

func HandleAddChirp(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		internal.RespondWithError(w, 401)
		return
	}
	decoder := json.NewDecoder(r.Body)
	chirp := internal.ChirpDto{}
	err := decoder.Decode(&chirp)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}

	if err := validateChirp(chirp); err != nil {
		internal.RespondWithError(w, 400)
		return
	}

	clean := cleanMessage(chirp.Body)
	newChirp, err := services.CreateChirp(userId, clean)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	internal.RespondWithJSON(w, 201, newChirp)
}

func HandleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := services.GetChirps()
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	internal.RespondWithJSON(w, 200, chirps)
}

func HandleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIdString := r.PathValue("chirpId")
	chirpId, err := strconv.Atoi(chirpIdString)
	if err != nil || chirpId < 0 {
		internal.RespondWithError(w, 404)
		return
	}
	chirp, err := services.GetChirpById(chirpId)
	if err != nil {
		if errors.Is(err, internal.ChirpNotFound) {
			internal.RespondWithError(w, 404)
			return
		}
		internal.RespondWithError(w, 500)
		return
	}
	internal.RespondWithJSON(w, 200, chirp)
}

func HandleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIdString := r.PathValue("chirpId")
	chirpId, err := strconv.Atoi(chirpIdString)
	if err != nil || chirpId < 0 {
		internal.RespondWithError(w, 404)
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		internal.RespondWithError(w, 401)
		return
	}
	_, err = services.DeleteChirpById(userId)
	if err != nil {
		internal.RespondWithError(w, 401)
		return
	}

	internal.RespondWithStatus(w, 204)
}

func validateChirp(chirp internal.ChirpDto) error {
	if len(chirp.Body) > 140 {
		return errors.New("Chirp is too long")
	}
	return nil
}

func cleanMessage(message string) string {
	words := strings.Fields(message)
	edited := false
	for i, word := range words {
		word = strings.ToLower(word)
		if !badWordTree.Exists(word) {
			continue
		}
		edited = true
		words[i] = "****"
	}
	if !edited {
		return message
	}
	return strings.Join(words, " ")
}
