package main

import (
	"fmt"
	"net/http"

	"github.com/JoshElias/chirpy/internal/handlers"
	"github.com/JoshElias/chirpy/internal/middleware"
	"github.com/joho/godotenv"
)

// listen on some endpoints, do some stuff
func main() {
	godotenv.Load()

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/app/", middleware.MetricsInc(http.StripPrefix("/app", fileServer)))

	mux.HandleFunc("GET /api/healthz", handlers.HealthHandler)
	mux.HandleFunc("/api/reset", handlers.ResetHandler)
	mux.HandleFunc("/admin/metrics", handlers.HandleMetricsAdmin)
	mux.HandleFunc("GET /api/chirps", handlers.HandleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}", handlers.HandleGetChirp)
	mux.HandleFunc("POST /api/chirps", handlers.HandleAddChirp)
	mux.HandleFunc("POST /api/users", handlers.HandleAddUser)
	mux.HandleFunc("POST /api/login", handlers.HandleLogin)
	mux.HandleFunc("PUT /api/users", handlers.HandleUpdateUser)

	fmt.Println("server listening on localhost:8080")
	if err := http.ListenAndServe("localhost:8080", mux); err == nil {
		fmt.Println("error starting server")
	}
}
