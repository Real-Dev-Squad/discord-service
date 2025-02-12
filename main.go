package main

import (
	"github.com/Real-Dev-Squad/discord-service/commands/register"
	config "github.com/Real-Dev-Squad/discord-service/config"
	queue "github.com/Real-Dev-Squad/discord-service/queue"
	"github.com/Real-Dev-Squad/discord-service/routes"
	"github.com/sirupsen/logrus"
)

func main() {
	register.SetupRegister()
	logrus.Info("Starting server on port " + config.AppConfig.Port)
	queue.GetQueueInstance()
	routes.Listen(":" + config.AppConfig.Port)
}
