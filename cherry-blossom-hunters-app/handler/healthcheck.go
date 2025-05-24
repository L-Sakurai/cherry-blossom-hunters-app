package handler

import (
	"encoding/json"
	"net/http"

	"cherry-blossom-hunters-app/logger"
	"cherry-blossom-hunters-app/service"
)

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	logger.Logging("/healthcheck accessed")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	res, err := service.FetchEvents()
	if err != nil {
		logger.Logging("eventDitect error: " + err.Error())
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Logging("Failed to encode JSON: " + err.Error())
		http.Error(w, `{"error": "Encoding Failed"}`, http.StatusInternalServerError)
	}
}
