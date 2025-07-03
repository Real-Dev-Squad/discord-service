package service

import (
	"fmt"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
)

func (s *CommandService) HelloService(response http.ResponseWriter, request *http.Request) {
	messageResponse := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Hey there <@%s>! Congratulations, you just executed your first slash command", s.discordMessage.Member.User.ID),
		},
	}
	utils.WriteJSONResponse(response, http.StatusOK, messageResponse)
}
