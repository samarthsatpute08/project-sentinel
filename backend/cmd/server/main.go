package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/samarthsatpute08/project-sentinel/internal/circuitbreaker"
	"github.com/samarthsatpute08/project-sentinel/internal/proxy"
	"github.com/samarthsatpute08/project-sentinel/internal/telemetry"
)

func main() {
	primaryURL := getEnv("PRIMARY_URL", "http://primary-api:8081")
	fallbackURL := getEnv("FALLBACK_URL", "http://fallback-api:8082")
	port := getEnv("PORT", "8080")

	// Create the telemetry hub (manages all WebSocket clients)
	hub := telemetry.NewHub()

	// Create circuit breaker with sensible defaults
	// 5 failures -> trip, 2 successes in half-open -> close, 10s timeout before probing
	breaker := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: 5,
		SuccessThreshold: 2,
		Timeout:          10 * time.Second,
	})

	breaker.OnStateChange = func(from, to circuitbreaker.State) {
		log.Printf("circuit breaker: %s -> %s", from, to)
	}

	// Create the router
	router, err := proxy.New(primaryURL, fallbackURL, breaker, hub)
	if err != nil {
		log.Fatalf("failed to create router: %v", err)
	}

	mux := http.NewServeMux()

	// All traffic proxied through our circuit-breaker router
	mux.Handle("/api/", router)

	// WebSocket endpoint for the React dashboard
	mux.Handle("/ws", hub.Handler())

	// Health check for Docker
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Serve React build (for production)
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("sentinel listening on :%s | primary=%s fallback=%s", port, primaryURL, fallbackURL)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}