package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	service "github.com/Real-Dev-Squad/discord-service/service"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/julienschmidt/httprouter"
)

func HomeHandler(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	payload, err := io.ReadAll(request.Body)
	if err != nil || len(payload) == 0 {
		utils.Errors.NewBadRequestError(response, "Invalid Request Payload")
		return
	}
	var message dtos.DiscordMessage
	err = json.Unmarshal(payload, &message)
	if err != nil {
		utils.Errors.NewInternalError(response)
		return
	}
	switch message.Type {

	case discordgo.InteractionPing:
		payload := map[string]interface{}{"type": uint8(discordgo.InteractionResponsePong)}
		utils.Success.NewDiscordResponse(response, "Pong", payload)
		return

	case discordgo.InteractionApplicationCommand:
		service.MainService(&message)(response, request)
		return

	default:
		response.WriteHeader(http.StatusOK)
	}
}
