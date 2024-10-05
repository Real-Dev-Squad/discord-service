package middleware

import (
	"encoding/hex"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/bwmarrin/discordgo"
	"github.com/julienschmidt/httprouter"
)

func VerifyCommand(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		publicKeyBytes, err := hex.DecodeString(config.AppConfig.DISCORD_PUBLIC_KEY)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		result := discordgo.VerifyInteraction(r, publicKeyBytes)
		if !result {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(w, r, ps)
	}
}
