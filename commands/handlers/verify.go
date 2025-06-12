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

var DISCORD_AVATAR_BASE_URL = "https://cdn.discordapp.com/avatars"
var VERIFICATION_STRING = "Please verify your discord account by clicking the link below 👇"
var VERIFICATION_SUBSTRING = "By granting authorization, you agree to permit us to manage your server nickname displayed ONLY in the Real Dev Squad server and to sync your joining data with your user account on our platform."

func (CS *CommandHandler) verify() error {
	metaData := CS.discordMessage.MetaData
	applicationId := metaData["applicationId"]
	
	uniqueToken := &utils.UniqueToken{}
	token, err := uniqueToken.GenerateUniqueToken()
	if err != nil {
		return err
	}

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.AppConfig.BOT_PRIVATE_KEY))
	if err != nil {
		return err
	}
	
	authToken := &utils.AuthToken{}
	authTokenString, err := authToken.GenerateAuthToken(jwt.SigningMethodRS256, jwt.MapClaims{
		"expiry": time.Now().Add(time.Second * 2).Unix(),
		"name": "Discord Service",
	}, rsaPrivateKey)
	if err != nil {
		return err
	}
	
	logrus.Infof("Auth token: %s generated at %v", authTokenString, time.Now().Unix())

	baseUrl:= fmt.Sprintf("%s/external-accounts", config.AppConfig.RDS_BASE_API_URL)
	requestBody:= map[string]any{
		"type": "discord",
		"token": token,
		"attributes": map[string]any{
			"discordId": CS.discordMessage.UserID,
			"userAvatar": fmt.Sprintf("%s/%s/%s.jpg", DISCORD_AVATAR_BASE_URL, CS.discordMessage.UserID, CS.discordMessage.MetaData["userAvatarHash"]),
			"userName": CS.discordMessage.MetaData["userName"],
			"discriminator": CS.discordMessage.MetaData["discriminator"],
			"discordJoinedAt": CS.discordMessage.MetaData["discordJoinedAt"],
			"expiry": time.Now().Add(time.Second * 2).Unix(),
		},
	}
	jsonBody, err := utils.Json.ToJson(requestBody)
	if err != nil {
		return err
	}
	
	request, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		return err
	}
	
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authTokenString))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-service-name", "Discord Service")
	
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	message:= ""
	if response.StatusCode == 201 || response.StatusCode == 200 {
		verificationSiteURL := "";
		if CS.discordMessage.MetaData["dev"] == "true" {
			verificationSiteURL = config.AppConfig.MAIN_SITE_URL;
			message = fmt.Sprintf("%s\n%s/discord?dev=true&token=%s\n%s", VERIFICATION_STRING, verificationSiteURL, token, VERIFICATION_SUBSTRING)
		}
		if metaData["dev"] == "false" {
			verificationSiteURL = config.AppConfig.VERIFICATION_SITE_URL;
			message = fmt.Sprintf("%s\n%s/discord?token=%s\n%s", VERIFICATION_STRING, verificationSiteURL, token, VERIFICATION_SUBSTRING)
		}
	}

	session, err := CreateSession()
	if err != nil {
		return err
	}

	webhookEdit := &discordgo.WebhookEdit{
		Content: &message,
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
