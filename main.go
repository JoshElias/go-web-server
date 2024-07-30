package main

import (
	"fmt"
	"net/http"
	"sync"
)

type apiConfig struct {
	mu             sync.Mutex
	fileserverHits int
}

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.mu.Lock()
		defer a.mu.Unlock()
		a.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	// mux.Handle("/", http.FileServer(http.Dir("./index.html")))
	// mux.HandleFunc("/app/*", func(w http.ResponseWriter, r *http.Request) {
	// 	r.URL.Path
	// 	http.ServeFile(w, r, "./index.html")
	// w.WriteHeader(http.StatusNotFound)
	// w.Write([]byte("404 - Not Found"))
	// })
	// mux.HandleFunc("/assets/logo.png", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./public/2CofkLc.png")
	// 	// w.WriteHeader(http.StatusNotFound)
	// w.Write([]byte("404 - Not Found"))
	// })
	apiConfig := apiConfig{}
	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	// http.Handle("/", fileServer)

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /api/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		apiConfig.mu.Lock()
		hits := apiConfig.fileserverHits
		apiConfig.mu.Unlock()
		w.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
	})

	mux.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		apiConfig.mu.Lock()
		apiConfig.fileserverHits = 0
		apiConfig.mu.Unlock()
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("server listening on localhost:8080")
	if err := http.ListenAndServe("localhost:8080", mux); err == nil {
		fmt.Println("error starting server")
	}

}
