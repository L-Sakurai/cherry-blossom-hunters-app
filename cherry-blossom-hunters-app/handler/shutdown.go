package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"cherry-blossom-hunters-app/logger"
)

func ShutdownHandler(done chan bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Logging("/shutdown accessed - shutting down server")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Shutdown request successfully requested.",
			"status":  "ok",
		})

		go func() {
			time.Sleep(1 * time.Second)
			done <- true
		}()
	}
}
