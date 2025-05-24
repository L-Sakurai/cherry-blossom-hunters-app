package service

import (
	"encoding/json"
	"os/exec"
)

type ScheduleEvent struct {
	Title       string `json:"title"`
	Level       string `json:"level"`
	Period      string `json:"period"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

func FetchEvents() ([]ScheduleEvent, error) {
	cmd := exec.Command("python3", "./script/event-scraper.py")
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
