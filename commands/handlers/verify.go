package handlers

import (
	"bytes"
	"fmt"
	"net/http"
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
	
	uniqueToken := &utils.UniqueToken{}
	token, err := uniqueToken.GenerateUniqueToken()
	if err != nil {
		return fmt.Errorf("error generating unique token: %v", err)
	}

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.AppConfig.BOT_PRIVATE_KEY))
	if err != nil {
		return fmt.Errorf("error parsing private key string to rsa private key: %v", err)
	}
	
	authToken := &utils.AuthToken{}
	authTokenString, err := authToken.GenerateAuthToken(jwt.SigningMethodRS256, jwt.MapClaims{
		"expiry": time.Now().Add(time.Second * 2).Unix(),
		"name": DISCORD_SERVICE,
	}, rsaPrivateKey)
	if err != nil {
		return fmt.Errorf("error generating auth token: %v", err)
	}

	baseUrl:= fmt.Sprintf("%s/external-accounts", config.AppConfig.RDS_BASE_API_URL)
	requestBody:= map[string]any{
		"type": "discord",
		"token": token,
		"attributes": map[string]any{
			"discordId": CS.discordMessage.UserID,
			"userAvatar": fmt.Sprintf("%s/%s/%s.jpg", DISCORD_AVATAR_BASE_URL, CS.discordMessage.UserID, metaData["userAvatarHash"]),
			"userName": metaData["userName"],
			"discriminator": metaData["discriminator"],
			"discordJoinedAt": metaData["discordJoinedAt"],
			"expiry": time.Now().Add(time.Second * 2).Unix(),
		},
	}
	
	jsonBody, err := utils.Json.ToJson(requestBody)
	if err != nil {
		return fmt.Errorf("error parsing request body in json string: %v", err)
	}
	
	request, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}
	
	request.Header.Set(HEADERS.AUTHORIZATION, fmt.Sprintf("Bearer %s", authTokenString))
	request.Header.Set(HEADERS.CONTENT_TYPE, "application/json")
	request.Header.Set(HEADERS.SERVICE, DISCORD_SERVICE)
	
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("error sending request to RDS Backend API: %v", err)
	}

	message:= "Something went wrong while generating verification link"
	if response.StatusCode == 201 || response.StatusCode == 200 {
		verificationSiteURL := ""
		if metaData["dev"] == "true" {
			verificationSiteURL = config.AppConfig.MAIN_SITE_URL;
			message = fmt.Sprintf("%s\n%s/discord?dev=true&token=%s\n%s", VERIFICATION_STRING, verificationSiteURL, token, VERIFICATION_SUBSTRING)
		}else if metaData["dev"] == "false" {
			verificationSiteURL = config.AppConfig.VERIFICATION_SITE_URL;
			message = fmt.Sprintf("%s\n%s/discord?token=%s\n%s", VERIFICATION_STRING, verificationSiteURL, token, VERIFICATION_SUBSTRING)
		}
	}

	session, err := CreateSession()
	if err != nil {
		return fmt.Errorf("error creating session: %v", err)
	}

	webhookEdit := &discordgo.WebhookEdit{
		Content: &message,
	}

	if _, err := session.WebhookMessageEdit(applicationId, metaData["token"], "@original", webhookEdit); err != nil {
		logrus.Errorf("Error editing original message for application %v", err)
		return fmt.Errorf("error editing original message for application %v", err)
	}

	if err := session.Close(); err != nil {
		logrus.Errorf("Error closing session: %v", err)
		return fmt.Errorf("error closing session: %v", err)
	}

	return nil
}
