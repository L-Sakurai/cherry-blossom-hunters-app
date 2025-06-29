package service

import (
	"context"
	"encoding/json"
	"os/exec"
	"time"
	"fmt"
	"strings"
	"cherry-blossom-hunters-app/utils/jsonUtil"
	"cherry-blossom-hunters-app/appConfig"
)

type ScheduleEvent struct {
	Title       string `json:"title"`
	Level       string `json:"level"`
	Period      string `json:"period"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

type EventFetchResponse struct {
	Events []ScheduleEvent
	Error  error
}

type EventService struct {
	scriptPath     string
	pythonCommand  string
	outputFilePath string
	defaultTimeout time.Duration
}

// NewEventService - Constructor for EventService
func NewEventService(appConfig *appConfig.Config) *EventService {
	config := appConfig
	return &EventService{
		scriptPath:     config.Event.ScriptPath,
		pythonCommand:  config.Event.PythonCommand,
		outputFilePath: config.Event.OutputFilePath,
		defaultTimeout: config.Event.DefaultTimeout,
	}
}

// FetchEvents - Fetch events with default timeout (maintains backward compatibility)
func (s *EventService) FetchEvents() ([]ScheduleEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.defaultTimeout)
	defer cancel()

	return s.FetchEventsWithContext(ctx)
}

// FetchEventsWithContext - Context-aware event fetching
func (s *EventService) FetchEventsWithContext(ctx context.Context) ([]ScheduleEvent, error) {
	cmd := exec.CommandContext(ctx, s.pythonCommand, s.scriptPath)
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("script execution failed: %w", err)
	}
	
	// Clean and validate JSON output
	cleanedOutput := s.cleanJSONOutput(output)
	
	// Log the cleaned output for debugging
	fmt.Printf("Cleaned JSON output: %s\n", string(cleanedOutput))
	
	// Persist the cleaned output
	if err := jsonUtil.PersistJSONToHash256(cleanedOutput, s.outputFilePath); err != nil {
		fmt.Printf("Warning: Failed to persist JSON to hash file: %v\n", err)
	}

	var res []ScheduleEvent
	if err := json.Unmarshal(cleanedOutput, &res); err != nil {
		return nil, fmt.Errorf("JSON parsing failed: %w, output: %s", err, string(cleanedOutput))
	}
	
	return res, nil
}

// cleanJSONOutput - Clean and validate JSON output from script
func (s *EventService) cleanJSONOutput(output []byte) []byte {
	outputStr := string(output)
	
	// Remove any leading/trailing whitespace
	outputStr = strings.TrimSpace(outputStr)
	
	// Find the start of JSON (should start with '[' or '{')
	startIdx := -1
	for i, char := range outputStr {
		if char == '[' || char == '{' {
			startIdx = i
			break
		}
	}
	
	if startIdx == -1 {
		// No JSON found, return original
		return output
	}
	
	// Find the end of JSON by counting brackets
	if startIdx < len(outputStr) {
		cleaned := outputStr[startIdx:]
		
		// Basic validation - ensure it starts and ends properly
		if strings.HasPrefix(cleaned, "[") && strings.HasSuffix(cleaned, "]") {
			return []byte(cleaned)
		} else if strings.HasPrefix(cleaned, "{") && strings.HasSuffix(cleaned, "}") {
			return []byte(cleaned)
		}
	}
	
	return output
}

// executeScript - Internal method for script execution (for future extensions)
func (s *EventService) executeScript(ctx context.Context) ([]byte, error) {
	cmd := exec.CommandContext(ctx, s.pythonCommand, s.scriptPath)
	return cmd.Output()
}

// processOutput - Internal method for output processing (for future extensions)
func (s *EventService) processOutput(output []byte) error {
	fmt.Println(string(output))
	return jsonUtil.PersistJSONToHash256(output, s.outputFilePath)
}

// parseEvents - Internal method for event parsing (for future extensions)
func (s *EventService) parseEvents(output []byte) ([]ScheduleEvent, error) {
	// Clean the output first
	cleanedOutput := s.cleanJSONOutput(output)
	
	// Validate JSON structure
	if !s.isValidJSON(cleanedOutput) {
		return nil, fmt.Errorf("invalid JSON structure: %s", string(cleanedOutput))
	}
	
	var res []ScheduleEvent
	if err := json.Unmarshal(cleanedOutput, &res); err != nil {
		return nil, fmt.Errorf("JSON unmarshaling failed: %w", err)
	}
	return res, nil
}

// isValidJSON - Check if the output is valid JSON
func (s *EventService) isValidJSON(data []byte) bool {
	var js interface{}
	return json.Unmarshal(data, &js) == nil
}

// validateEventData - Validate individual event data
func (s *EventService) validateEventData(events []ScheduleEvent) error {
	for i, event := range events {
		if event.Title == "" {
			return fmt.Errorf("event at index %d has empty title", i)
		}
		if event.Level == "" {
			return fmt.Errorf("event '%s' at index %d has empty level", event.Title, i)
		}
		// Add more validation as needed
	}
	return nil
}

// Legacy functions for backward compatibility
// These may be deprecated in the future

// FetchEvents - Legacy function (global)
func FetchEvents() ([]ScheduleEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return FetchEventsWithContext(ctx)
}

// FetchEventsWithContext - Legacy function (global)
func FetchEventsWithContext(ctx context.Context) ([]ScheduleEvent, error) {
	cmd := exec.CommandContext(ctx, "python3", "./script/event-scraper.py")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	fmt.Println(string(output))
	jsonUtil.PersistJSONToHash256(output, "./sha256_output.txt")

	var res []ScheduleEvent
	if err := json.Unmarshal(output, &res); err != nil {
		return nil, err
	}
	
	return res, nil
}