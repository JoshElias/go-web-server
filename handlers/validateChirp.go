package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func HandleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body  string `json:"body"`
		Error string `json:"error"`
	}

	type returnObj struct {
		Error string `json:"error,omitempty"`
		Valid bool   `json:"valid"`
	}

	errorResponse := func(err error, code int) {
		res := returnObj{
			Error: fmt.Sprintf("%v", err),
		}
		data, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(data)
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorResponse(err, 500)
		return
	}

	if len(params.Body) > 140 {
		err := errors.New("Chirp is too long")
		errorResponse(err, 400)
		return
	}

	res := returnObj{
		Valid: true,
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
