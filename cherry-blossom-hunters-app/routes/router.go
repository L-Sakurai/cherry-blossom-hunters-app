package routes

import (
	"net/http"
	"cherry-blossom-hunters-app/handler"
	"cherry-blossom-hunters-app/appConfig"
)

func SetupRoutes(httpShutdown chan bool, appConfig *appConfig.Config) *http.ServeMux {
	mux := http.NewServeMux()
	
	// Member compliance handler with Discord notification
	memberComplianceHandler := handler.NewMemberServiceHandler(appConfig)
	mux.HandleFunc("/api/function/compliance/check", memberComplianceHandler.CheckCompliance)
	
	// Event check handler with Discord notification
	eventServiceHandler := handler.NewEventServiceHandler(appConfig)

	mux.HandleFunc("/api/function/event/check", eventServiceHandler.CheckEvents)
	
	// System handlers
	mux.HandleFunc("/api/systems/healthcheck", handler.HealthcheckHandler)
	mux.HandleFunc("api/systems/shutdown", handler.ShutdownHandler(httpShutdown))
	
	return mux
}