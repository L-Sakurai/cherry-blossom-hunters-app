package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cherry-blossom-hunters-app/logger"
	"cherry-blossom-hunters-app/service"
	"cherry-blossom-hunters-app/appConfig"
	"cherry-blossom-hunters-app/notify"
)

type EventServiceHandler struct {
	eventService *service.EventService
	notifier     *notify.Notifier
}

func NewEventServiceHandler(cfg *appConfig.Config) *EventServiceHandler {
	return &EventServiceHandler{
		eventService: service.NewEventService(cfg),
		notifier:     notify.NewDiscordNotifier(cfg.Notify.Webhooks["event"]),
	}
}

func (h *EventServiceHandler) CheckEvents(w http.ResponseWriter, r *http.Request) {
	logger.Logging("info", "/events accessed")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.eventService.FetchEventsWithContext(ctx)
	if err != nil {
		logger.Logging("error", "eventDetect error: %v", err)

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
			"error":     errorMessage,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	response := map[string]interface{}{
		"events":    res,
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	// Discord通知の送信
	h.sendEventNotification(res)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logging("error", "Failed to encode JSON: %v", err)
		http.Error(w, `{"error": "Encoding Failed"}`, http.StatusInternalServerError)
	}
}

func (h *EventServiceHandler) sendEventNotification(events []service.ScheduleEvent) {
	if len(events) == 0 {
		return
	}
	
	payload := &notify.CustomPayload{
		Title:   "📅 イベント検知機能",
		Message: fmt.Sprintf("新しいイベントが %d 件検出されました", len(events)),
		Fields: map[string]interface{}{
			"検出されたイベント": events,
			"イベント数":      len(events),
		},
	}

	if err := h.notifier.Send(payload); err != nil {
		logger.Logging("error", "イベント通知送信失敗: %v", err)
	} else {
		logger.Logging("info", "イベント通知に成功しました")
	}
}

// 後方互換性のため、既存のEventHandler関数も残す
// Global EventService instance for EventHandler
var eventService *service.EventService

// InitEventService - Initialize EventService (call at application startup)
func InitEventService(config *appConfig.Config) {
	eventService = service.NewEventService(config)	
}

func EventHandler(w http.ResponseWriter, r *http.Request) {
	logger.Logging("info", "/events accessed")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Fetch events using EventService
	// If eventService is not initialized, use legacy function (backward compatibility)
	var res []service.ScheduleEvent
	var err error
	
	if eventService != nil {
		res, err = eventService.FetchEventsWithContext(ctx)
	} else {
		// Fallback: use legacy function
		res, err = service.FetchEventsWithContext(ctx)
	}
	
	if err != nil {
		logger.Logging("error", "eventDetect error: %v", err)

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
			"error":     errorMessage,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	response := map[string]interface{}{
		"events":    res,
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logging("error", "Failed to encode JSON: %v", err)
		http.Error(w, `{"error": "Encoding Failed"}`, http.StatusInternalServerError)
	}
}