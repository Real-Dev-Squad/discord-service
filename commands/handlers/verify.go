package handlers

import (
	"fmt"
	"time"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func (CS *CommandHandler) verify() error {
	metaData := CS.discordMessage.MetaData
	applicationId := metaData["applicationId"]
	
	_, err := utils.TokenHelper.GenerateUniqueToken()
	if err != nil {
		return err
	}

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.AppConfig.BOT_PRIVATE_KEY))
	if err != nil {
		return err
	}

	authToken, err:= utils.TokenHelper.GenerateAuthToken(jwt.SigningMethodRS256, jwt.MapClaims{
		"expiry": time.Now().Add(time.Second * 2).Unix(),
		"name": "Discord Service",
	}, rsaPrivateKey)
	if err != nil {
		return err
	}
	
	logrus.Infof("Auth token: %s generated at %v", authToken, time.Now().Unix())

	session, err := CreateSession()
	if err != nil {
		return err
	}
	
	// Todo: This is a temporary solution to check that message is being sent correctly.
	content := fmt.Sprintf("Hey there <@%s>! Congratulations, you just executed your first slash command", CS.discordMessage.UserID)
	
	webhookEdit := &discordgo.WebhookEdit{
		Content: &content,
	}
	
	if _, err := session.WebhookMessageEdit(applicationId, metaData["token"], "@original", webhookEdit); err != nil {
		logrus.Errorf("Error editing webhook message for application %s: %v", applicationId, err)
		return err
	}

	if err := session.Close(); err != nil {
		logrus.Errorf("Error closing session: %v", err)
		return err
	}
	return nil
}
