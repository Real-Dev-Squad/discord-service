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
	rw.Header().Set("content-type", "application/json;charset=UTF-8")
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Error(err)
	}
	var message dtos.DiscordMessage
	err = json.Unmarshal(payload, &message)
	if err != nil {
		logrus.Error(err)
	}
	switch message.Type {
	case discordgo.InteractionApplicationCommand:
		resp := map[string]uint8{"type": uint8(discordgo.InteractionResponseDeferredChannelMessageWithSource)}
		err = json.NewEncoder(rw).Encode(resp)
		if err != nil {
			logrus.Error(err)
		}
		return
	default:
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
}
