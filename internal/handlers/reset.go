package handlers

import (
	"github.com/JoshElias/go-web-server/internal"
	"net/http"
)

func ResetHandler(w http.ResponseWriter, r *http.Request) {
	m := internal.GetMetrics()
	m.Mu.Lock()
	m.FileserverHits = 0
	m.Mu.Unlock()
	w.WriteHeader(http.StatusOK)
}
