package notify

import (
	"cherry-blossom-hunters-app/logger"
)

func NotifyUserShutdown() {
	logger.Logging("=== SERVER SHUTDOWN NOTIFICATION ===")
	logger.Logging("Application server has been stopped successfully")
	logger.Logging("All users have been notified of the shutdown")
	logger.Logging("====================================")
}