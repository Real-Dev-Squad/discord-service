package middleware

import (
	"encoding/hex"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/bwmarrin/discordgo"
	"github.com/julienschmidt/httprouter"
)

func VerifyCommand(next httprouter.Handle) httprouter.Handle {
	return func(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
		publicKeyBytes, err := hex.DecodeString(config.AppConfig.DISCORD_PUBLIC_KEY)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		result := discordgo.VerifyInteraction(request, publicKeyBytes)
		if !result {
			response.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(response, request, params)
	}
}
