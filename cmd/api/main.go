package main

import (
	"context"
	"log"
	netHttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-checker/internal/repository"
	"url-checker/internal/service"
	"url-checker/internal/transport/http"
)

func main() {
	repo, err := repository.NewRepository("./storage/tasks.json")
	if err != nil {
		log.Fatalf("Failed to init storage: %v", err)
	}
	checkerService := service.NewCheckerService(repo)
	r := http.NewRouter(checkerService)

	srv := &netHttp.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != netHttp.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced shutdown: %v", err)
	}
	if err := checkerService.Shutdown(5 * time.Second); err != nil {
		log.Printf("Service shutdown error: %v", err)
	}

	log.Println("Server stopped gracefully")
}
