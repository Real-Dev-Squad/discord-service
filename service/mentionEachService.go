package service

import (
	"fmt"
	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/queue"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *CommandService) MentionEachService(response http.ResponseWriter, request *http.Request) {
	logrus.Info("Mention-each command received")

	if s.discordMessage == nil || s.discordMessage.Data == nil ||
		s.discordMessage.Member == nil || s.discordMessage.Member.User == nil {
		errorMsg := "Invalid request data"
		logrus.Error(errorMsg)
		sendErrorResponse(response, errorMsg)
		return
	}

	var roleID, message string
	var dev, devTitle bool

	roleOption := findOption(s.discordMessage.Data.Options, "role")
	if roleOption == nil {
		errorMsg := "Role is required for mention-each command"
		logrus.Error(errorMsg)
		sendErrorResponse(response, errorMsg)
		return
	}

	roleID, ok := roleOption.Value.(string)
	if !ok {
		errorMsg := "Invalid role format(not a String)"
		logrus.Errorf("%s: Expected string, got %T", errorMsg, roleOption.Value)
		sendErrorResponse(response, errorMsg)
		return
	}

	if roleID == "" {
		errorMsg := "Invalid role format (empty ID)"
		logrus.Error(errorMsg)
		sendErrorResponse(response, errorMsg)
		return
	}

	messageOption := findOption(s.discordMessage.Data.Options, "message")
	if messageOption != nil && messageOption.Value != nil {
		message, _ = messageOption.Value.(string)
	}

	devOption := findOption(s.discordMessage.Data.Options, "dev")
	if devOption != nil && devOption.Value != nil {
		dev, _ = devOption.Value.(bool)
	}

	devTitleOption := findOption(s.discordMessage.Data.Options, "dev_title")
	if devTitleOption != nil && devTitleOption.Value != nil {
		devTitle, _ = devTitleOption.Value.(bool)
	}

	logrus.Infof("Mention-each options: role=%s, message=%s, dev=%v, devTitle=%v", roleID, message, dev, devTitle)

	var responseContent string
	if devTitle {
		responseContent = fmt.Sprintf("Fetching users with the <@&%s>", roleID)
	} else if dev {
		responseContent = fmt.Sprintf("Sending individual mentions to users with the <@&%s>", roleID)
	} else {
		responseContent = fmt.Sprintf("Mentioning all users with the <@&%s>", roleID)
	}

	metaData := map[string]string{
		"role_id":    roleID,
		"message":    message,
		"dev":        fmt.Sprintf("%v", dev),
		"dev_title":  fmt.Sprintf("%v", devTitle),
		"guild_id":   s.discordMessage.Data.GuildId,
		"channel_id": s.discordMessage.ChannelId,
	}

	dataPacket := &dtos.DataPacket{
		UserID:      s.discordMessage.Member.User.ID,
		CommandName: utils.CommandNames.MentionEach,
		MetaData:    metaData,
	}

	bytePacket, err := dataPacket.ToByte()
	if err != nil {
		errorMsg := "Failed to process your request"
		logrus.Errorf("Failed to marshal message: %v", err)
		sendErrorResponse(response, errorMsg)
		return
	}

	if err := queue.SendMessage(bytePacket); err != nil {
		errorMsg := "Failed to process your request"
		logrus.Errorf("Failed to enqueue message: %v", err)
		sendErrorResponse(response, errorMsg)
		return
	}

	messageResponse := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseContent,
			Flags:   0,
		},
	}

	utils.Success.NewDiscordResponse(response, "Success", messageResponse)

}

func findOption(options []*discordgo.ApplicationCommandInteractionDataOption, name string) *discordgo.ApplicationCommandInteractionDataOption {
	for _, opt := range options {
		if opt.Name == name {
			return opt
		}
	}

	return nil
}

func sendErrorResponse(response http.ResponseWriter, errorMessage string) {
	messageResponse := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: errorMessage,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
	utils.Success.NewDiscordResponse(response, "Error", messageResponse)
}
