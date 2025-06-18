package notify

import (
	"context"
	"time"

	"cherry-blossom-hunters-app/logger"
)

// NotifyUserShutdown is the existing function (no changes)
func NotifyUserShutdown() {
	// Implement existing notification logic here
	logger.Logging("Sending shutdown notifications to users...")
	
	// Example: actual notification processing (email, Slack, webhook, etc.)
	// This part depends on the original code
	time.Sleep(200 * time.Millisecond) // Simulation of notification processing
	
	logger.Logging("User shutdown notifications sent successfully")
}

// NotifyUserShutdownWithContext is the context-aware version
func NotifyUserShutdownWithContext(ctx context.Context) error {
	logger.Logging("Sending shutdown notifications to users...")
	
	// Check for context cancellation while processing notifications
	select {
	case <-ctx.Done():
		logger.Logging("Notification canceled due to context timeout", logger.Warn)
		return ctx.Err()
	default:
		// Actual notification processing
		// Example: external API calls, database updates, etc.
		time.Sleep(200 * time.Millisecond) // Simulation of notification processing
	}
	
	logger.Logging("User shutdown notifications sent successfully")
	return nil
}

// NotifyUserShutdownWithTimeout is notification with timeout
func NotifyUserShutdownWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	return NotifyUserShutdownWithContext(ctx)
}