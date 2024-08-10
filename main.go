package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type apiConfig struct {
	mu             sync.Mutex
	fileserverHits int
}

// Middleware
func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.mu.Lock()
		defer a.mu.Unlock()
		a.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

// listen on some endpoints, do some stuff
func main() {
	mux := http.NewServeMux()
	apiConfig := apiConfig{}

	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		apiConfig.mu.Lock()
		apiConfig.fileserverHits = 0
		apiConfig.mu.Unlock()
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		var template = `<html>
<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>
`
		w.Header().Set("Content-Type", "text/html")
		apiConfig.mu.Lock()
		hits := apiConfig.fileserverHits
		apiConfig.mu.Unlock()
		w.Write([]byte(fmt.Sprintf(template, hits)))
	})

	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
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
	})

	fmt.Println("server listening on localhost:8080")
	if err := http.ListenAndServe("localhost:8080", mux); err == nil {
		fmt.Println("error starting server")
	}

}
