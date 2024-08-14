package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/JoshElias/chirpy/internal"
)

type ChirpDTO struct {
	Body string `json:"body"`
}

type ChirpEnitity struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

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
	chirp := internal.Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}

	if err := validateChirp(chirp); err != nil {
		internal.RespondWithError(w, 400)
	}

	chirp.Body = cleanMessage(chirp.Body)
	// save in database
	// return new chirp entity
	internal.RespondWithJSON(w, 201, chirp)
}

func HandleGetChirps(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func validateChirp(chirp internal.Chirp) error {
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
