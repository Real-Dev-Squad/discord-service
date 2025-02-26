package service

import (
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/utils"
)

type CommandService struct {
	discordMessage *dtos.DiscordMessage
}

func (s *CommandService) MainService(response http.ResponseWriter, request *http.Request) {
	switch s.discordMessage.Data.Name {
	case utils.CommandNames.Hello:
		s.HelloService(response, request)
		return
	case utils.CommandNames.Listening:
		s.ListeningService(response, request)
		return
	case utils.CommandNames.Verify:
		s.VerifyService(response, request)
		return
	default:
		utils.Errors.NewBadRequestError(response, "Invalid Command")
		return
	}
}
