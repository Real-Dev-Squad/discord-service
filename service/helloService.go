package service

import (
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
)

func (s *CommandService) HelloService(response http.ResponseWriter, request *http.Request) {
	messageResponse := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: utils.ResponseGenerator.HelloResponse(s.discordMessage.Member.User.ID),
		},
	}
	utils.Success.NewDiscordResponse(response, "Success", messageResponse)
}
