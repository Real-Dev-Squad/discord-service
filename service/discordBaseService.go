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
	if err != nil || len(payload) == 0 {
		utils.Errors.NewBadRequestError(response, "Invalid Request Payload")
		return
	}
	defer request.Body.Close()
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
		MainService(&message)(response, request)
		return

	default:
		response.WriteHeader(http.StatusOK)
	}
}
