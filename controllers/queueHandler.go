package controllers

import (
	"io"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/commands/handlers"
	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func QueueHandler(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(response, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	handler := handlers.MainHandler(body)
	if handler != nil {
		if err := utils.ExponentialBackoffRetry(config.AppConfig.MAX_RETRIES, handler); err != nil {
			logrus.Errorf("Failed to process command after %d attempts: %s", config.AppConfig.MAX_RETRIES, err)
		}
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

}
