package handlers

import (
	"errors"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/models"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/sirupsen/logrus"
)

type CommandHandler struct {
	discordMessage *dtos.DataPacket
}

var CS = CommandHandler{}

func MainHandler(dataPacket []byte) func() error {
	packetData := &dtos.DataPacket{}
	err := utils.FromByte(dataPacket, packetData)
	if err != nil {
		logrus.Errorf("Failed to unmarshal data send by queue: %v", err)
		return nil
	}
	CS.discordMessage = packetData
	switch packetData.CommandName {
	case utils.CommandNames.Listening:
		return CS.listeningHandler
	default:
		logrus.Warn("Invalid Command Received: ", packetData.CommandName)
		return nil
	}
}

func UpdateNickName(userId string, newNickName string, sessionWrapper models.SessionInterface) error {
	if len(newNickName) > 32 {
		logrus.Error("Must be 32 or fewer in length.")
		return errors.New("Must be 32 or fewer in length.")
	}
	if err := sessionWrapper.GuildMemberNickname(userId, newNickName); err != nil {
		logrus.Errorf("Failed to update nickname: %v", err)
		return err
	}
	return nil
}
