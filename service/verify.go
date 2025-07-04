package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/errors"
	"github.com/Real-Dev-Squad/discord-service/queue"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (s *CommandService) Verify(response http.ResponseWriter, request *http.Request) {
	dev := "false"
	if len(s.discordMessage.Data.Options) == 1 {
		if val, ok := s.discordMessage.Data.Options[0].Value.(bool); ok && val{
			dev = "true"
		}
	}
	
	message:= &dtos.DataPacket{
		UserID:      s.discordMessage.Member.User.ID,
		CommandName: utils.CommandNames.Verify,
		MetaData: map[string]string{
			"userAvatarHash":  s.discordMessage.Member.Avatar,
			"userName":        s.discordMessage.Member.User.Username,
			"discriminator":   s.discordMessage.Member.User.Discriminator,
			"discordJoinedAt": s.discordMessage.Member.JoinedAt.Format(time.RFC3339),
			"dev":             dev,
			"channelId":       s.discordMessage.ChannelId,
			"token":           s.discordMessage.Token,
			"applicationId":   s.discordMessage.ApplicationId,
		},
	}

	messageBytes, err:= json.Marshal(message)
	if err != nil {
		logrus.Errorf("Failed to convert data packet to json bytes: %v", err)
		errors.HandleError(response, err)
		return
	}

	if err := queue.SendMessage(messageBytes); err != nil {
		logrus.Errorf("Failed to send data packet to queue: %v", err)
		errors.HandleError(response, err)
		return
	}

	res:= &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Your request is being processed.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}

	utils.WriteJSONResponse(response, http.StatusOK, res)
}
