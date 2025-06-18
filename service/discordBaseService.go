package service

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/errors"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
)

func DiscordBaseService(response http.ResponseWriter, request *http.Request) {
	payload, err := io.ReadAll(request.Body)
	if err != nil || len(payload) == 0 {
		errors.HandleError(response, errors.NewBadRequest("Invalid Request Payload", err))
		return
	}
	defer request.Body.Close()
	
	
	var message dtos.DiscordMessage
	if err = json.Unmarshal(payload, &message); err != nil {
		errors.HandleError(response, errors.NewInternalServerError("Internal Server Error", err))
		return
	}

	switch message.Type {

	case discordgo.InteractionPing:
		payload := dtos.DiscordMessage{
			Type:    discordgo.InteractionPing,
		}
		utils.WriteJSONResponse(response, http.StatusOK, map[string]any{
			"message": "Pong",
			"data": payload,
		})
		return

	case discordgo.InteractionApplicationCommand:
		MainService(&message)(response, request)
		return

	default:
		errors.HandleError(response, errors.NewBadRequest("Invalid Request Payload", err))
		return
	}
}
