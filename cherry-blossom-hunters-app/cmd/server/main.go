package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	// "fmt"

	"cherry-blossom-hunters-app/routes"
	"cherry-blossom-hunters-app/logger"
	"cherry-blossom-hunters-app/notify"
	"cherry-blossom-hunters-app/appConfig"
	// "cherry-blossom-hunters-app/service"
)

func main() {
	config := appConfig.GetConfig()
	logger.SetUp()
	// Channel for graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Channel for HTTP shutdown
	httpShutdown := make(chan bool, 1)

	// HTTP server configuration
	server := &http.Server{
		Addr:    ":8080",
		Handler: routes.SetupRoutes(httpShutdown, config),
	}

	// Goroutine for server startup
	go func() {
		logger.Logging("Server listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logging("Server error: "+err.Error(), logger.Error)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	select {
	case sig := <-shutdown:
		logger.Logging("Received signal: " + sig.String() + ". Initiating graceful shutdown...")
	case <-httpShutdown:
		logger.Logging("HTTP shutdown request received. Initiating graceful shutdown...")
	}

	// Execute graceful shutdown
	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	// Set shutdown timeout (30 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// Execute external notifications concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Logging("Sending shutdown notifications...")
		notify.NotifyUserShutdown()
	}()

	// HTTP server shutdown
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Logging("Shutting down HTTP server...")
		if err := server.Shutdown(ctx); err != nil {
			logger.Logging("Server shutdown error: "+err.Error(), logger.Error)
		} else {
			logger.Logging("HTTP server shutdown complete")
		}
	}()

	// Wait for all shutdown processes to complete
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Logging("Graceful shutdown completed successfully")
	case <-ctx.Done():
		logger.Logging("Shutdown timeout exceeded, forcing exit", logger.Warn)
	}

	// Final cleanup
	time.Sleep(100 * time.Millisecond)
	logger.Logging("Application exiting...")
}