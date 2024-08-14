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
	}

	clean := cleanMessage(chirp.Body)
	dbPath := filepath.Join(getCwd(), "database.json")
	conn, err := internal.NewDbConnection(dbPath)
	if err != nil {
		internal.RespondWithError(w, 500)
	}
	newChirp, err := conn.CreateChirp(clean)
	if err != nil {
		internal.RespondWithError(w, 500)
	}
	// save in database
	// return new chirp entity
	internal.RespondWithJSON(w, 201, newChirp)
}

func HandleGetChirps(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func getCwd() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
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
