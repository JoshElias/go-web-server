package handlers

import (
	"github.com/JoshElias/chirpy/config"
	"net/http"
)

func ResetHandler(w http.ResponseWriter, r *http.Request) {
	c := config.GetConfig()
	c.Mu.Lock()
	c.FileserverHits = 0
	c.Mu.Unlock()
	w.WriteHeader(http.StatusOK)
}
