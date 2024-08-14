package internal

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}
