package register

import (
	constants "github.com/Real-Dev-Squad/discord-service/commands"
	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/models"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var NewDiscord = discordgo.New

func SetupRegister() {
	session, err := NewDiscord("Bot " + config.AppConfig.BOT_TOKEN)
	if err != nil {
		logrus.Panic("Cannot create a new Discord session ", err)
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logrus.Info("Logged in as: ", session.State.User.Username, session.State.User.Discriminator)
	})
	sessionWrapper := &models.SessionWrapper{Session: session}
	RegisterCommands(sessionWrapper)
}

var RegisterCommands = func(openSession models.SessionInterface) {
	err := openSession.Open()
	defer openSession.Close()

	if err != nil {
		logrus.Panic("Cannot open the session ", err)
	}

	for _, v := range constants.Commands {
		_, err := openSession.ApplicationCommandCreate(openSession.GetUserId(), config.AppConfig.GUILD_ID, v)
		if err != nil {
			logrus.Panic("Cannot create ", v.Name, "command: ", err)
		}
	}
	logrus.Info("Successfully registered commands")
}
