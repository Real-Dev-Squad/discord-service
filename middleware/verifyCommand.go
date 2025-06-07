package middleware

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

var VerifyInteraction = discordgo.VerifyInteraction

func VerifyCommand(next httprouter.Handle) httprouter.Handle {
	return func(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
		response.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// Read the request body
		body, err := io.ReadAll(request.Body)
		if err != nil {
			logrus.Errorf("Failed to read request body: %v", err)
			utils.Errors.NewInternalError(response)
			return
		}
		// Restore the request body for the next handler
		request.Body.Close()
		request.Body = io.NopCloser(bytes.NewReader(body))

		// Verify the Discord signature
		publicKeyBytes, err := hex.DecodeString(config.AppConfig.DISCORD_PUBLIC_KEY)
		if err != nil {
			logrus.Errorf("Failed to decode Discord public key: %v", err)
			utils.Errors.NewInternalError(response)
			return
		}

		result := VerifyInteraction(request, publicKeyBytes)
		if !result {
			logrus.Warn("Unauthorized request attempt")
			utils.Errors.NewUnauthorisedError(response)
			return
		}

		next(response, request, params)
	}
}
