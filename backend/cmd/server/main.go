// Command server runs the calculator HTTP API.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	calculatorbackend "calculator-backend"
	"calculator-backend/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", handler.Calculate)
	mux.HandleFunc("/healthz", healthz)

	rootHandler := calculatorbackend.WithCORS(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           rootHandler,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	// Run the server in a background goroutine so the main goroutine is
	// free to wait for a shutdown signal below.
	go func() {
		log.Printf("server listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	// Block until we receive an interrupt or termination signal.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutdown signal received, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		return
	}

	log.Println("server stopped")
}

// healthz reports that the service is running. It always returns 200 with a
// small JSON body; it does not check downstream dependencies because this
// service has none.
func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
