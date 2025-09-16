package service

import (
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/utils"
)

type CommandService struct {
	discordMessage *dtos.DiscordMessage
}

var CS = CommandService{}

func MainService(discordMessage *dtos.DiscordMessage) func(response http.ResponseWriter, request *http.Request) {
	CS.discordMessage = discordMessage
	switch discordMessage.Data.Name {
	case utils.CommandNames.Hello:
		return CS.HelloService
	case utils.CommandNames.Listening:
		return CS.ListeningService
	case utils.CommandNames.Verify:
		return CS.Verify
	default:
		return func(response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(http.StatusOK)
		}
	}
}
