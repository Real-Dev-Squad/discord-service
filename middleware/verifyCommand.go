package middleware

import (
	"encoding/hex"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/errors"
	"github.com/bwmarrin/discordgo"
	"github.com/julienschmidt/httprouter"
)

var VerifyInteraction = discordgo.VerifyInteraction

func VerifyCommand(next httprouter.Handle) httprouter.Handle {
	return func(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
		response.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		publicKeyBytes, err := hex.DecodeString(config.AppConfig.DISCORD_PUBLIC_KEY)
		if err != nil {
			errors.HandleError(response, errors.NewInternalServerError("Internal Server Error", err))
			return
		}
		
		result := VerifyInteraction(request, publicKeyBytes)
		if !result {
			errors.HandleError(response, errors.NewUnauthorized("Unauthorized Access", nil))
			return
		}
		
		next(response, request, params)
	}
}
