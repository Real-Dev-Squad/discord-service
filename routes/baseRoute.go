package routes

import (
	"github.com/Real-Dev-Squad/discord-service/controllers"
	"github.com/Real-Dev-Squad/discord-service/middleware"
	"github.com/gorilla/mux"
)

func SetupBaseRoutes(router *mux.Router) {
	router.HandleFunc("/", middleware.VerifyCommand(controllers.HomeHandler)).Methods("POST")
	router.HandleFunc("/health", controllers.HealthCheckHandler).Methods("GET")
}
