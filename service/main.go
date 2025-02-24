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

// MainService returns an HTTP handler function that dispatches a Discord command to the appropriate service method.
// It assigns the provided Discord message to the service state and uses the command name within the message data to
// select the corresponding method. For recognized commands (Hello, Listening, Verify), it returns the respective service;
// otherwise, it returns a default handler that writes a 200 OK status to the response.
func MainService(discordMessage *dtos.DiscordMessage) func(response http.ResponseWriter, request *http.Request) {
	CS.discordMessage = discordMessage
	switch discordMessage.Data.Name {
	case utils.CommandNames.Hello:
		return CS.HelloService
	case utils.CommandNames.Listening:
		return CS.ListeningService
	case utils.CommandNames.Verify:
		return CS.VerifyService
	default:
		return func(response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(http.StatusOK)
		}
	}
}
