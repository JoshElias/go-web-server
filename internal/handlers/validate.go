package handlers

import (
	"encoding/json"
	"net/http"
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

func HandleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body  string `json:"body"`
		Error string `json:"error"`
	}

	type returnObj struct {
		CleanedBody string `json:"cleaned_body,omitempty"`
	}

	errorResponse := func(code int) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorResponse(500)
		return
	}

	if len(params.Body) > 140 {
		errorResponse(400)
		return
	}

	res := returnObj{
		CleanedBody: cleanChirp(params.Body),
	}
	data, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

func cleanChirp(chirp string) string {
	words := strings.Fields(chirp)
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
		return chirp
	}
	return strings.Join(words, " ")
}
