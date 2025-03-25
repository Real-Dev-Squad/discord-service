package handlers

import (
	"errors"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/models"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type CommandHandler struct {
	discordMessage *dtos.DataPacket
}

var CS = CommandHandler{}

func MainHandler(dataPacket []byte) func() error {
	packetData := &dtos.DataPacket{}
	err := packetData.FromByte(dataPacket)
	if err != nil {
		logrus.Errorf("Failed to unmarshal data send by queue: %v", err)
		return nil
	}
	CS.discordMessage = packetData
	switch packetData.CommandName {
	case utils.CommandNames.Listening:
		return CS.listeningHandler
	case utils.CommandNames.MentionEach:
		return CS.mentionEachHandler
	default:
		logrus.Warn("Invalid Command Received: ", packetData.CommandName)
		return nil
	}
}

type DiscordSession struct {
	session *discordgo.Session
}

var NewDiscord = discordgo.New
var CreateSession = func() (*discordgo.Session, error) {
	session, err := NewDiscord("Bot " + config.AppConfig.BOT_TOKEN)
	if err != nil {
		logrus.Errorf("Cannot create a new Discord session: %v", err)
		return nil, err
	}
	openSession := &models.SessionWrapper{Session: session}
	err = openSession.Open()
	if err != nil {
		logrus.Errorf("Cannot open the session: %v", err)
		return nil, err
	}
	return session, nil
}

func UpdateNickName(userId string, newNickName string) error {
	if len(newNickName) > 32 {
		logrus.Error("Must be 32 or fewer in length.")
		return errors.New("Must be 32 or fewer in length.")
	}
	session, err := CreateSession()
	if err != nil {
		return err
	}
	err = session.GuildMemberNickname(config.AppConfig.GUILD_ID, userId, newNickName)
	if err != nil {
		logrus.Errorf("Cannot update nickname: %v", err)
		return nil
	}
	return nil
}

func (c *CommandHandler) mentionEachHandler() error {
	logrus.Info("mentionEach command received")
	return nil
}
