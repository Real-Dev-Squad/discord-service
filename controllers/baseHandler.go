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

func HomeHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Error(err)
		http.Error(rw, "Failed to read request body", http.StatusBadRequest)
		return
	}
	var message dtos.DiscordMessage
	err = json.Unmarshal(payload, &message)
	if err != nil {
		logrus.Error(err)
		http.Error(rw, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	switch message.Type {
	case discordgo.InteractionPing:
		resp := map[string]uint8{"type": uint8(discordgo.InteractionResponsePong)}
		err = json.NewEncoder(rw).Encode(resp)
		if err != nil {
			logrus.Error(err)
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		rw.WriteHeader(http.StatusOK)
	}
}
