package main

import (
	"fmt"
	"net/http"

	"github.com/JoshElias/chirpy/config"
	"github.com/JoshElias/chirpy/handlers"
	"github.com/JoshElias/chirpy/middleware"
)

// listen on some endpoints, do some stuff
func main() {
	mux := http.NewServeMux()
	apiConfig := config.ApiConfig{}
	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/app/", middleware.MetricsInc(http.StripPrefix("/app", fileServer), &apiConfig))

	mux.HandleFunc("GET /api/healthz", handlers.HealthHandler)
	mux.Handle("/api/reset", handlers.ResetHandler(&apiConfig))
	// 	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
	// 		var template = `<html>
	// <body>
	//     <h1>Welcome, Chirpy Admin</h1>
	//     <p>Chirpy has been visited %d times!</p>
	// </body>
	//
	// </html>
	// `
	// 		w.Header().Set("Content-Type", "text/html")
	// 		apiConfig.mu.Lock()
	// 		hits := apiConfig.fileserverHits
	// 		apiConfig.mu.Unlock()
	// 		w.Write([]byte(fmt.Sprintf(template, hits)))
	// 	})

	// mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
	// 	type parameters struct {
	// 		Body  string `json:"body"`
	// 		Error string `json:"error"`
	// 	}
	//
	// 	type returnObj struct {
	// 		Error string `json:"error,omitempty"`
	// 		Valid bool   `json:"valid"`
	// 	}
	//
	// 	errorResponse := func(err error, code int) {
	// 		res := returnObj{
	// 			Error: fmt.Sprintf("%v", err),
	// 		}
	// 		data, err := json.Marshal(res)
	// 		if err != nil {
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			return
	// 		}
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(code)
	// 		w.Write(data)
	// 	}
	//
	// 	decoder := json.NewDecoder(r.Body)
	// 	params := parameters{}
	// 	err := decoder.Decode(&params)
	// 	if err != nil {
	// 		errorResponse(err, 500)
	// 		return
	// 	}
	//
	// 	if len(params.Body) > 140 {
	// 		err := errors.New("Chirp is too long")
	// 		errorResponse(err, 400)
	// 		return
	// 	}
	//
	// 	res := returnObj{
	// 		Valid: true,
	// 	}
	// 	data, err := json.Marshal(res)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(200)
	// 	w.Write(data)
	// })
	// 	type parameters struct {
	// 		Body  string `json:"body"`
	// 		Error string `json:"error"`
	// 	}
	//
	// 	type returnObj struct {
	// 		Error string `json:"error,omitempty"`
	// 		Valid bool   `json:"valid"`
	// 	}
	//
	// 	errorResponse := func(err error, code int) {
	// 		res := returnObj{
	// 			Error: fmt.Sprintf("%v", err),
	// 		}
	// 		data, err := json.Marshal(res)
	// 		if err != nil {
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			return
	// 		}
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(code)
	// 		w.Write(data)
	// 	}
	//
	// 	decoder := json.NewDecoder(r.Body)
	// 	params := parameters{}
	// 	err := decoder.Decode(&params)
	// 	if err != nil {
	// 		errorResponse(err, 500)
	// 		return
	// 	}
	//
	// 	if len(params.Body) > 140 {
	// 		err := errors.New("Chirp is too long")
	// 		errorResponse(err, 400)
	// 		return
	// 	}
	//
	// 	res := returnObj{
	// 		Valid: true,
	// 	}
	// 	data, err := json.Marshal(res)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(200)
	// 	w.Write(data)
	// })
	//
	fmt.Println("server listening on localhost:8080")
	if err := http.ListenAndServe("localhost:8080", mux); err == nil {
		fmt.Println("error starting server")
	}

}
