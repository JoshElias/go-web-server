package middleware

import (
	"github.com/JoshElias/chirpy/internal"
	"net/http"
)

func MetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := internal.GetMetrics()
		m.Mu.Lock()
		defer m.Mu.Unlock()
		m.FileserverHits++
		next.ServeHTTP(w, r)
	})
}
