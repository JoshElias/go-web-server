package main

import (
	"fmt"
	"net/http"

	"github.com/JoshElias/chirpy/internal/handlers"
	"github.com/JoshElias/chirpy/internal/middleware"
)

// listen on some endpoints, do some stuff
func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/app/", middleware.MetricsInc(http.StripPrefix("/app", fileServer)))

	mux.HandleFunc("GET /api/healthz", handlers.HealthHandler)
	mux.HandleFunc("/api/reset", handlers.ResetHandler)
	mux.HandleFunc("/admin/metrics", handlers.HandleMetricsAdmin)
	mux.HandleFunc("/api/validate_chirp", handlers.HandleValidateChirp)

	fmt.Println("server listening on localhost:8080")
	if err := http.ListenAndServe("localhost:8080", mux); err == nil {
		fmt.Println("error starting server")
	}
}
