package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/JoshElias/chirpy/internal"
	"github.com/JoshElias/chirpy/internal/services"
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
	newChirp, err := services.CreateChirp(clean)
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
	conn, err := internal.GetTestDbConnection()
	if err != nil {
		internal.RespondWithError(w, 500)
	}
	db, err := conn.LoadDb()
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}

	chirp, exists := db.Chirps[chirpId]
	if !exists {
		internal.RespondWithError(w, 404)
		return
	}
	internal.RespondWithJSON(w, 200, chirp)
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
