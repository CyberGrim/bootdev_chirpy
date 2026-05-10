package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(".")))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	serveErr := srv.ListenAndServe()
	if serveErr != nil {
		log.Fatal(serveErr)
	}
}
