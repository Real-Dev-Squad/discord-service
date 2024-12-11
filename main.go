package main

import (
	config "github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/routes"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Starting server on port " + config.AppConfig.Port)
	routes.Listen(":" + config.AppConfig.Port)
}
