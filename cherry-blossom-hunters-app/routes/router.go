package routes

import "net/http"

import "cherry-blossom-hunters-app/handler"


func SetupRoutes(httpShutdown chan bool) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck", handler.HealthcheckHandler)
	mux.HandleFunc("/shutdown", handler.ShutdownHandler(httpShutdown))
	return mux
}
