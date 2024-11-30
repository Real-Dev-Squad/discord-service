package routes

import (
	"github.com/Real-Dev-Squad/discord-service/controllers"
	"github.com/Real-Dev-Squad/discord-service/middleware"
	"github.com/julienschmidt/httprouter"
)

func SetupBaseRoutes(router *httprouter.Router) {
	router.POST("/", middleware.VerifyCommand(controllers.HomeHandler))

	// Add the health check route
	router.GET("/health", controllers.HealthCheckHandler)
}
