package service

import (
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/utils"
)

type CommandService struct {
	discordMessage *dtos.DiscordMessage
}

func (service *CommandService) HandleMessage(response http.ResponseWriter, request *http.Request) {
	switch service.discordMessage.Data.Name {
	case utils.CommandNames.Hello:
		service.Hello(response, request)
		return
	case utils.CommandNames.Listening:
		service.Listening(response, request)
		return
	case utils.CommandNames.Verify:
		service.Verify(response, request)
		return
	default:
		utils.Errors.NewBadRequestError(response, "Invalid Command")
		return
	}
}
