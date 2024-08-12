package handlers

import (
	"github.com/JoshElias/chirpy/config"
	"net/http"
)

func ResetHandler(a *config.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.Mu.Lock()
		a.FileserverHits = 0
		a.Mu.Unlock()
		w.WriteHeader(http.StatusOK)
	})
}
