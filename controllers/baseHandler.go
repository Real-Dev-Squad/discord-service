package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/julienschmidt/httprouter"
)

func HomeHandler(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	payload, err := io.ReadAll(request.Body)
	if err != nil {
		utils.Errors.NewBadRequestError(response, "Failed to read request body")
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
	default:
		response.WriteHeader(http.StatusOK)
	}
}
