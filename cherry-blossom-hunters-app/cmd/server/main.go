package main

import (
	"net/http"
	"os"
	"time"

	"cherry-blossom-hunters-app/handler"
	"cherry-blossom-hunters-app/logger"
	"cherry-blossom-hunters-app/notify"
)

func main() {
	logger.SetUp()

	done := make(chan bool)

	http.HandleFunc("/healthcheck", handler.HealthcheckHandler)
	http.HandleFunc("/shutdown", handler.ShutdownHandler(done))

	logger.Logging("Server listening on :8080")
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			logger.Logging("Server error: " + err.Error())
			os.Exit(1)
		}
	}()

	if <-done {
		logger.Logging("Shutdown signal received. Initiating graceful shutdown...")
		notify.NotifyUserShutdown()
		time.Sleep(500 * time.Millisecond)
		logger.Logging("Shutdown complete. Exiting...")
	}
}
