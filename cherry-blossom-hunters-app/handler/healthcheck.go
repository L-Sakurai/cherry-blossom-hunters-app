package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"cherry-blossom-hunters-app/logger"
	"cherry-blossom-hunters-app/service"
)

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	logger.Logging("/healthcheck accessed")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// リクエストのタイムアウトを設定（5秒）
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := service.FetchEventsWithContext(ctx)
	if err != nil {
		logger.Logging("eventDetect error: "+err.Error(), logger.Error)
		
		// エラーの種類に応じて適切なステータスコードを返す
		var statusCode int
		var errorMessage string
		
		if ctx.Err() == context.DeadlineExceeded {
			statusCode = http.StatusRequestTimeout
			errorMessage = "Request timeout"
		} else if ctx.Err() == context.Canceled {
			statusCode = http.StatusRequestTimeout
			errorMessage = "Request canceled"
		} else {
			statusCode = http.StatusInternalServerError
			errorMessage = "Internal server error"
		}
		
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{
			"error": errorMessage,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	// 成功レスポンス
	response := map[string]interface{}{
		"events": res,
		"status": "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logging("Failed to encode JSON: "+err.Error(), logger.Error)
		http.Error(w, `{"error": "Encoding Failed"}`, http.StatusInternalServerError)
	}
}
