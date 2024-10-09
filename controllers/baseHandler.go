package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/bwmarrin/discordgo"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func HomeHandler(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	payload, err := io.ReadAll(request.Body)
	if err != nil {
		logrus.Error(err)
		http.Error(response, "Failed to read request body", http.StatusBadRequest)
		return
	}
	var message dtos.DiscordMessage
	err = json.Unmarshal(payload, &message)
	if err != nil {
		logrus.Error(err)
		http.Error(response, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	switch message.Type {
	case discordgo.InteractionPing:
		resp := map[string]uint8{"type": uint8(discordgo.InteractionResponsePong)}
		err = json.NewEncoder(response).Encode(resp)
		if err != nil {
			logrus.Error(err)
			http.Error(response, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		response.WriteHeader(http.StatusOK)
	}
}
