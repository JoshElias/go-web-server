package main

import (
	"fmt"
	"net/http"

	"github.com/JoshElias/go-web-server/internal/handlers"
	"github.com/JoshElias/go-web-server/internal/middleware"
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
	mux.Handle(
		"POST /api/chirps",
		middleware.Auth(
			http.HandlerFunc(handlers.HandleAddChirp),
		),
	)
	mux.HandleFunc("POST /api/users", handlers.HandleAddUser)
	mux.HandleFunc("POST /api/login", handlers.HandleLogin)
	mux.Handle(
		"PUT /api/users",
		middleware.Auth(
			http.HandlerFunc(handlers.HandleUpdateUser),
		),
	)
	mux.HandleFunc("POST /api/refresh", handlers.HandleRefreshToken)
	mux.HandleFunc("POST /api/revoke", handlers.HandleTokenRevoke)

	fmt.Println("server listening on localhost:8080")
	if err := http.ListenAndServe("localhost:8080", mux); err == nil {
		fmt.Println("error starting server")
	}
}
