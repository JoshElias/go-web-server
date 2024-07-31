package main

import (
	"encoding/json"
	"fmt"
	"log"
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

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding body: %s", err)
			res := returnObj{
				Error: fmt.Sprintf("%v", err),
			}
			data, err := json.Marshal(res)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write(data)
			return
		}

		if len(params.Body) > 140 {
			res := returnObj{
				Error: "Chirp is too long",
			}
			data, err := json.Marshal(res)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Println("sending error data")
			fmt.Println(data)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write(data)
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
