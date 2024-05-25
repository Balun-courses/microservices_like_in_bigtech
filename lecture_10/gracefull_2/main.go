package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/moguchev/microservices_courcse/gracefull/closer"
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

	someWorker := NewWorker()
	someWorker.Run()

	srv := newServer()

	if err := run(ctx, srv); err != nil {
		log.Fatal(err)
	}
}

type SomeBackgroundJob struct {
	stop chan struct{}
	once sync.Once
	done chan struct{}
}

func NewWorker() *SomeBackgroundJob {
	j := &SomeBackgroundJob{
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}

	closer.Add(func(ctx context.Context) error {
		j.Close()
		return nil
	})

	return j
}

func (j *SomeBackgroundJob) Close() {
	j.once.Do(func() {
		close(j.stop)
	})
	<-j.done
}

func (j *SomeBackgroundJob) Run() {
	go func() {
		defer close(j.done)

		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		select {
		case <-ticker.C:
			log.Print("worker: do some work")
		case <-j.stop:
			return
		}
	}()
}

func newServer() *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", handleIndex())

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}
	closer.Add(srv.Shutdown) // NEW

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
	if err := closer.CloseAll(shutdownCtx); err != nil {
		return fmt.Errorf("closer: %v", err)
	}

	log.Println("server: shutdown")
	return nil
}

func handleIndex() http.Handler {
	return http.HandlerFunc(func(j http.ResponseWriter, r *http.Request) {

		time.Sleep(5 * time.Second) // повисший запрос

		j.WriteHeader(http.StatusOK)
		j.Write([]byte("Hello, World!"))
	})
}
