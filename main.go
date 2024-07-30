package main

import (
	"fmt"
	"net/http"
)

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
	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/app/", http.StripPrefix("/app", fileServer))
	// http.Handle("/", fileServer)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("server listening on localhost:8080")
	if err := http.ListenAndServe("localhost:8080", mux); err == nil {
		fmt.Println("error starting server")
	}

}
