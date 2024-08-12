package middleware

import (
	"github.com/JoshElias/chirpy/config"
	"net/http"
)

func MetricsInc(next http.Handler, a *config.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.Mu.Lock()
		defer a.Mu.Unlock()
		a.FileserverHits++
		next.ServeHTTP(w, r)
	})
}
