package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/JoshElias/chirpy/internal"
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
	wd, err := os.Getwd()
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	dbPath := filepath.Join(wd, "database.json")
	conn, err := internal.NewDbConnection(dbPath)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	newChirp, err := conn.CreateChirp(clean)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	internal.RespondWithJSON(w, 201, newChirp)
}

func HandleGetChirps(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	dbPath := filepath.Join(wd, "database.json")
	conn, err := internal.NewDbConnection(dbPath)
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	chirps, err := conn.GetChirps()
	if err != nil {
		internal.RespondWithError(w, 500)
		return
	}
	internal.RespondWithJSON(w, 200, chirps)
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
