package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"cherry-blossom-hunters-app/logger"
)

func ShutdownHandler(shutdown chan bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// POSTリクエストのみ受け付ける（セキュリティ向上）
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

		// レスポンスを送信
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Logging("Failed to encode shutdown response: "+err.Error(), logger.Error)
			return
		}

		// レスポンスのフラッシュを確実に行う
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		// 非同期でシャットダウンシグナルを送信
		go func() {
			// レスポンスが確実に送信されるまで少し待機
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