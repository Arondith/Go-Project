package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	gh "github.com/Arondith/Portfolio-project/internal/github"
)

func main() {
	client := gh.NewClient(os.Getenv("GITHUB_TOKEN"))

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/api/analyze", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		owner := strings.TrimSpace(r.URL.Query().Get("owner"))
		repo := strings.TrimSpace(r.URL.Query().Get("repo"))
		if owner == "" || repo == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "owner and repo query parameters are required"})
			return
		}

		analysis, err := client.Analyze(r.Context(), owner, repo)
		if err != nil {
			log.Printf("analyze %s/%s: %v", owner, repo, err)
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, analysis)
	})

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           requestLogger(mux),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("RepoRadar listening on http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(started).Round(time.Millisecond))
	})
}
