package service

import (
	"fmt"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
)

func (s *CommandService) VerifyService(response http.ResponseWriter, request *http.Request) {
	msg := "Work in progress for Verify command"
	messageResponse := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(msg),
			Flags:   64, // Ephemeral message flag
		},
	}
	utils.Success.NewDiscordResponse(response, "Success", messageResponse)
}
