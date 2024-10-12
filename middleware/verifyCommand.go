package middleware

import (
	"encoding/hex"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/julienschmidt/httprouter"
)

func VerifyCommand(next httprouter.Handle) httprouter.Handle {
	return func(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
		publicKeyBytes, err := hex.DecodeString(config.AppConfig.DISCORD_PUBLIC_KEY)
		if err != nil {
			utils.Errors.NewInternalError(response)
			return
		}
		result := discordgo.VerifyInteraction(request, publicKeyBytes)
		if !result {
			utils.Errors.NewUnauthorisedError(response)
			return
		}
		next(response, request, params)
	}
}
