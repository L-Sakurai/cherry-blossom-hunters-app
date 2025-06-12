package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"cherry-blossom-hunters-app/logger"
)

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	logger.Logging("/healthcheck accessed")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response := map[string]string{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logging("Failed to encode JSON: "+err.Error(), logger.Error)
		http.Error(w, `{"error": "Encoding Failed"}`, http.StatusInternalServerError)
	}
}
