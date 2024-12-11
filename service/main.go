package service

import (
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/dtos"
)

type CommandService struct {
	discordMessage *dtos.DiscordMessage
}

var CS = CommandService{}

func MainService(discordMessage *dtos.DiscordMessage) func(response http.ResponseWriter, request *http.Request) {
	CS.discordMessage = discordMessage
	switch discordMessage.Data.Name {
	case "hello":
		return CS.HelloService

	default:
		return func(response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(http.StatusOK)
		}
	}
}
