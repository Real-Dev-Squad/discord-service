package routes

import (
	"github.com/Real-Dev-Squad/discord-service/controllers"
	"github.com/Real-Dev-Squad/discord-service/middleware"
	"github.com/julienschmidt/httprouter"
)

// SetupBaseRoutes configures the base HTTP routes for the application. It registers a POST route at "/" that wraps the DiscordBaseHandler with the VerifyCommand middleware, a GET route at "/health" for health checks via HealthCheckHandler, and a POST route at "/queue" for queue operations managed by QueueHandler.
func SetupBaseRoutes(router *httprouter.Router) {
	router.POST("/", middleware.VerifyCommand(controllers.DiscordBaseHandler))
	router.GET("/health", controllers.HealthCheckHandler)
	router.POST("/queue", controllers.QueueHandler)
}
