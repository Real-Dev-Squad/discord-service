package service

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
)

func DiscordBaseService(response http.ResponseWriter, request *http.Request) {
	payload, err := io.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil || len(payload) == 0 {
		utils.Errors.NewBadRequestError(response, "Invalid Request Payload")
		return
	}
	var message dtos.DiscordMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		utils.Errors.NewInternalError(response)
		return
	}
	switch message.Type {

	case discordgo.InteractionPing:
		payload := map[string]interface{}{"type": uint8(discordgo.InteractionResponsePong)}
		utils.Success.NewDiscordResponse(response, "Pong", payload)
		return

	case discordgo.InteractionApplicationCommand:
		mainService := CommandService{discordMessage: message}
		mainService.HandleMessage(response, request)

		return

	default:
		response.WriteHeader(http.StatusNotFound)
	}
}
