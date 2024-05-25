package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const (
	listenAddr = "127.0.0.1:8080"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	srv := newServer()

	if err := run(ctx, srv); err != nil {
		log.Fatal(err)
	}
}

func newServer() *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", handleIndex())

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}
	return srv
}

func run(ctx context.Context, srv *http.Server) error {
	// Start HTTP server in a goroutine
	go func() {
		log.Println("server: start serve", listenAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server: %v", err)
		}
	}()

	// Wait until we receive a shutdown signal
	<-ctx.Done()

	log.Println("server: shutting down server gracefully")

	// Create a context with a 20-second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Attempt a graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server: shutdown: %w", err)
	}

	log.Println("server: shutdown")
	return nil
}

func handleIndex() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		time.Sleep(5 * time.Second) // повисший запрос

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})
}
