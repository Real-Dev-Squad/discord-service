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
			Content: "Hey there! Congratulations, you just executed your first slash command",
		},
	}
	utils.Success.NewDiscordResponse(response, "Success", messageResponse)
}
