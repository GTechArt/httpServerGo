package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	// atomic.Int32: safe to read/write from multiple goroutines
	// Use the “atomic” method to avoid a concurrency situation when multiple
	// requests increment the counter
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const apipath = "/api"
	const port = "8080"

	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	//mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", HandleReadiness)
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, req *http.Request) {
		s := fmt.Sprintf(`<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>`,
			apiCfg.fileserverHits.Load())
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(s))
	})
	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, req *http.Request) {
		apiCfg.fileserverHits.Store(0)
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func HandleReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		//log.Printf("counter = %d", cfg.fileserverHits.Load())
		next.ServeHTTP(w, req)
	})
}
