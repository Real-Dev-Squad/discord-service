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
	defer request.Body.Close()
	if err != nil || len(payload) == 0 {
		errors.HandleError(response, errors.NewBadRequest("Invalid Request Payload", err))
		return
	}

	var message dtos.DiscordMessage
	if err = json.Unmarshal(payload, &message); err != nil {
		errors.HandleError(response, err)
		return
	}

	switch message.Type {

	case discordgo.InteractionPing:
        utils.WriteJSONResponse(response, http.StatusOK, map[string]uint8{
         "type": uint8(discordgo.InteractionResponsePong),
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
