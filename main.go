package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	metricsLoad := cfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf("Hits: %d", metricsLoad)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	metricsLoad := cfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf("Hits: %d. Metrics Reset.", metricsLoad)))
}

func main() {
	mux := http.NewServeMux()
	api := &apiConfig{}

	mux.Handle("/app/", http.StripPrefix("/app", api.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/metrics", api.handlerMetrics)
	mux.HandleFunc("POST /api/reset", api.handlerReset)
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	serveErr := srv.ListenAndServe()
	if serveErr != nil {
		log.Fatal(serveErr)
	}
}
