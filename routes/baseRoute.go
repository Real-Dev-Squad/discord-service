package routes

import (
	"github.com/Real-Dev-Squad/discord-service/controllers"
	"github.com/Real-Dev-Squad/discord-service/middleware"
	"github.com/julienschmidt/httprouter"
)

func SetupBaseRoutes(router *httprouter.Router) {
	router.POST("/", middleware.VerifyCommand(controllers.HomeHandler))
	router.GET("/health", controllers.HealthCheckHandler)
	router.POST("/queue", controllers.QueueHandler)
}
