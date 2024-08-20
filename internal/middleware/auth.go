package middleware

import (
	"github.com/JoshElias/go-web-server/internal"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := internal.GetMetrics()
		m.Mu.Lock()
		defer m.Mu.Unlock()
		m.FileserverHits++
		next.ServeHTTP(w, r)
	})
}
