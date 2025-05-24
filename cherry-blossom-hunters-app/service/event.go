package service

import (
	"context"
	"encoding/json"
	"os/exec"
	"time"
)

type ScheduleEvent struct {
	Title       string `json:"title"`
	Level       string `json:"level"`
	Period      string `json:"period"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

// コンテキスト対応のFetchEvents関数
func FetchEvents() ([]ScheduleEvent, error) {
	// デフォルトで10秒のタイムアウトを設定
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return FetchEventsWithContext(ctx)
}

// コンテキストを受け取るバージョン
func FetchEventsWithContext(ctx context.Context) ([]ScheduleEvent, error) {
	cmd := exec.CommandContext(ctx, "python3", "./script/event-scraper.py")
	
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var res []ScheduleEvent
	if err := json.Unmarshal(output, &res); err != nil {
		return nil, err
	}
	
	return res, nil
}