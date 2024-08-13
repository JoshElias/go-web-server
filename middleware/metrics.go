package middleware

import (
	"github.com/JoshElias/chirpy/config"
	"net/http"
)

func MetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := config.GetConfig()
		c.Mu.Lock()
		defer c.Mu.Unlock()
		c.FileserverHits++
		next.ServeHTTP(w, r)
	})
}
