package controllers

import (
	"net/http"

	services "github.com/Real-Dev-Squad/discord-service/services"
	"github.com/julienschmidt/httprouter"
)

func DiscordBaseHandler(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	services.DiscordBaseService(response, request)
	return
}
