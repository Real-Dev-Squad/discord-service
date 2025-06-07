package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/commands/handlers"
	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func QueueHandler(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		logrus.Errorf("Failed to read request body: %v", err)
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(Response{
			Status:  "error",
			Message: "Failed to read request body",
		})
		return
	}

	handler := handlers.MainHandler(body)
	if handler == nil {
		logrus.Error("No handler found for the request")
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(Response{
			Status:  "error",
			Message: "Invalid request",
		})
		return
	}

	if err := utils.ExponentialBackoffRetry(config.AppConfig.MAX_RETRIES, handler); err != nil {
		logrus.Errorf("Failed to process command after %d attempts: %v", config.AppConfig.MAX_RETRIES, err)
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(Response{
			Status:  "error",
			Message: "Failed to process command",
		})
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(Response{
		Status: "success",
	})
}
