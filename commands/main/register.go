package main

import (
	constants "github.com/Real-Dev-Squad/discord-service/commands"
	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/models"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	session, err := utils.NewDiscordSession("Bot " + config.AppConfig.BOT_TOKEN)
	if err != nil {
		logrus.Panic("Cannot create a new Discord session", err)
	}

	sessionWrapper := &models.SessionWrapper{Session: session}
	RegisterCommands(sessionWrapper)
}

var RegisterCommands = func(openSession models.SessionInterface) {
	openSession.AddHandler()
	err := openSession.Open()
	if err != nil {
		logrus.Panic("Cannot open the session ", err)
	}

	for _, v := range constants.Commands {
		_, err := openSession.ApplicationCommandCreate(openSession.GetUerId(), config.AppConfig.GUILD_ID, v)
		if err != nil {
			logrus.Panic("Cannot create ", v.Name, "command: ", err)
		}
	}

	defer openSession.Close()
}
