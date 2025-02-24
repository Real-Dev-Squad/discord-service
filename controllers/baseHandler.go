package controllers

import (
	"net/http"

	service "github.com/Real-Dev-Squad/discord-service/service"
	"github.com/julienschmidt/httprouter"
)

// DiscordBaseHandler handles incoming HTTP requests for Discord by delegating the request processing to DiscordBaseService.
// It accepts an http.ResponseWriter, a *http.Request, and httprouter.Params (provided by the router but not utilized by this handler).
func DiscordBaseHandler(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	service.DiscordBaseService(response, request)
	return
}
