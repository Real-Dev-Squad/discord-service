package controllers

import (
	"net/http"

	service "github.com/Real-Dev-Squad/discord-service/service"
	"github.com/julienschmidt/httprouter"
)

func DiscordBaseHandler(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	service.DiscordBaseService(response, request)
	return
}
