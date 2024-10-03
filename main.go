package main

import (
	"net/http"

	config "github.com/Real-Dev-Squad/discord-service/config"
	utility "github.com/Real-Dev-Squad/discord-service/utils"
)

var logger = &utility.Logger{}

func main() {
	logger.Info("Starting server")
	err := http.ListenAndServe(":"+config.Config.Port, nil)
	if err != nil {
		logger.Error(err)
	}
}
