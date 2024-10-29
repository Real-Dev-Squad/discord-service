package main

import (
	constants "github.com/Real-Dev-Squad/discord-service/commands"
	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type sessionWrapper struct {
	session *discordgo.Session
}

func (s *sessionWrapper) open() error {
	return s.session.Open()
}

func (s *sessionWrapper) close() error {
	return s.session.Close()
}
func (s *sessionWrapper) applicationCommandCreate(applicationID, guildID string, command *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error) {
	return s.session.ApplicationCommandCreate(applicationID, guildID, command)
}

func (sw *sessionWrapper) getUerId() string {
	return sw.session.State.User.ID
}

type SessionInterface interface {
	open() error
	close() error
	applicationCommandCreate(applicationID, guildID string, command *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error)
	getUerId() string
}

var NewDiscord = discordgo.New

func main() {
	SetupConnection()
}

var SetupConnection = func() {
	var err error
	var session *discordgo.Session
	session, err = NewDiscord("Bot " + config.AppConfig.BOT_TOKEN)
	if err != nil {
		logrus.Error("Cannot create a new Discord session: ")
		panic(err)
	}
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logrus.Info("Logged in as: ", session.State.User.Username, session.State.User.Discriminator)
	})
	sessionWrapper := &sessionWrapper{session: session}

	RegisterCommands(sessionWrapper)
}

var RegisterCommands = func(openSession SessionInterface) {
	err := openSession.open()
	if err != nil {
		logrus.Error("Cannot open the session: ")
		panic(err)
	}
	for _, v := range constants.Commands {
		_, err := openSession.applicationCommandCreate(openSession.getUerId(), config.AppConfig.GUILD_ID, v)
		if err != nil {
			logrus.Error("Cannot create ", v.Name, "command: ", err)
			panic(err)
		}
	}
	defer openSession.close()
}
