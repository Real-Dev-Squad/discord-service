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
		logrus.Error("Cannot create a new Discord session: ")
		panic(err)
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logrus.Info("Logged in as: ", session.State.User.Username, session.State.User.Discriminator)
	})
	sessionWrapper := &models.SessionWrapper{Session: session}
	RegisterCommands(sessionWrapper)
}

var RegisterCommands = func(openSession models.SessionInterface) {
	err := openSession.Open()
	if err != nil {
		logrus.Error("Cannot open the session: ")
		panic(err)
	}

	for _, v := range constants.Commands {
		_, err := openSession.ApplicationCommandCreate(openSession.GetUerId(), config.AppConfig.GUILD_ID, v)
		if err != nil {
			logrus.Error("Cannot create ", v.Name, "command: ", err)
			panic(err)
		}
	}
	logrus.Info("Successfully registered commands")

	defer openSession.Close()
}
