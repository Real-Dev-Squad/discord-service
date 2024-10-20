package main

import (
	constants "github.com/Real-Dev-Squad/discord-service/commands"
	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

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
	err = session.Open()
	if err != nil {
		logrus.Error("Cannot open the session: ")
		panic(err)
	}
	logrus.Info("Adding commands...")
	for _, v := range constants.Commands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, config.AppConfig.GUILD_ID, v)
		if err != nil {
			logrus.Error("Cannot create ", v.Name, "command: ", err)
			panic(err)
		}
	}
	defer session.Close()
}
