package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"cherry-blossom-hunters-app/logger"
)

func ShutdownHandler(shutdown chan bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		logger.Logging("/shutdown accessed - initiating server shutdown")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		response := map[string]interface{}{
			"message":   "Shutdown request received successfully.",
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Logging("Failed to encode shutdown response: "+err.Error(), logger.Error)
			return
		}

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		go func() {
			time.Sleep(100 * time.Millisecond)
			
			select {
			case shutdown <- true:
				logger.Logging("Shutdown signal sent successfully")
			default:
				logger.Logging("Shutdown channel is full or closed", logger.Warn)
			}
		}()
	}
}