package main

import (
	"context"
	"log"
	"my-clean-architecture-template/config"
	"my-clean-architecture-template/internal/app"
	"my-clean-architecture-template/pkg/helper"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(3) // 3 goroutine

	go func() {
		defer wg.Done()
		app.RunWeb(ctx, cfg)
	}()

	go func() {
		defer wg.Done()
		app.RunWorker(ctx, cfg)
	}()

	go func() {
		defer wg.Done()
		helper.CleanUp(ctx)
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
	log.Println("Received interrupt signal. Shutting down...")

	cancel()

	wg.Wait()
	log.Println("Shutdown complete.")
}
