package models

import (
	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type SessionWrapper struct {
	Session *discordgo.Session
}

func (s *SessionWrapper) Open() error {
	return s.Session.Open()
}

func (s *SessionWrapper) Close() error {
	return s.Session.Close()
}

func (s *SessionWrapper) ApplicationCommandCreate(applicationID, guildID string, command *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error) {
	return s.Session.ApplicationCommandCreate(applicationID, guildID, command)
}

func (sw *SessionWrapper) GetUerId() string {
	return sw.Session.State.User.ID
}
func (sw *SessionWrapper) AddHandler() func() {
	return sw.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logrus.Info("Logged in as: ", sw.Session.State.User.Username, sw.Session.State.User.Discriminator)
	})
}
func (sw *SessionWrapper) GuildMemberNickname(userId string, newNickName string) error {
	return sw.Session.GuildMemberNickname(config.AppConfig.GUILD_ID, userId, newNickName)
}

type SessionInterface interface {
	AddHandler() func()
	Open() error
	Close() error
	ApplicationCommandCreate(applicationID, guildID string, command *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error)
	GetUerId() string
	GuildMemberNickname(userId string, newNickName string) error
}

var CreateSession = func() (*SessionWrapper, error) {
	session, err := utils.NewDiscordSession("Bot " + config.AppConfig.BOT_TOKEN)
	if err != nil {
		logrus.Errorf("Cannot create a new Discord session: %v", err)
		return nil, err
	}
	newSessionWrapper := &SessionWrapper{Session: session}
	return newSessionWrapper, nil
}
