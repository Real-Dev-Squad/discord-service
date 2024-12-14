package middleware

import (
	"encoding/hex"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
)

var VerifyInteraction = discordgo.VerifyInteraction

func VerifyCommand(next http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json; charset=UTF-8")
		publicKeyBytes, err := hex.DecodeString(config.AppConfig.DISCORD_PUBLIC_KEY)
		if err != nil {
			utils.Errors.NewInternalError(response)
			return
		}
		result := VerifyInteraction(request, publicKeyBytes)
		if !result {
			utils.Errors.NewUnauthorisedError(response)
			return
		}
		next(response, request)
	}
}
