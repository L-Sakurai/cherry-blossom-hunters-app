package routes

import "net/http"

import "cherry-blossom-hunters-app/handler"
import "cherry-blossom-hunters-app/appConfig"


func SetupRoutes(httpShutdown chan bool, appConfig *appConfig.Config) *http.ServeMux {
	mux := http.NewServeMux()
	memberComplianceHandler := handler.NewMemberServiceHandler(appConfig)
	mux.HandleFunc("/api/member/compliance/check", memberComplianceHandler.CheckCompliance)
	mux.HandleFunc("/healthcheck", handler.HealthcheckHandler)
	mux.HandleFunc("/shutdown", handler.ShutdownHandler(httpShutdown))
	
	return mux
}
