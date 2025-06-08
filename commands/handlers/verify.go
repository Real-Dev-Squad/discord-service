package handlers

import (
	"bytes"
<<<<<<< HEAD
	"encoding/json"
=======
>>>>>>> 4f1bed6 (fix(verify): add logic to make a http call with website-backend and content generation)
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
	
<<<<<<< HEAD
	token, err := utils.UniqueToken.GenerateUniqueToken()
=======
	token, err := utils.TokenHelper.GenerateUniqueToken()
>>>>>>> 4f1bed6 (fix(verify): add logic to make a http call with website-backend and content generation)
	if err != nil {
		return fmt.Errorf("error generating unique token: %v", err)
	}

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.AppConfig.BOT_PRIVATE_KEY))
	if err != nil {
		return fmt.Errorf("error parsing private key string to rsa private key: %v", err)
	}
	
	authTokenString, err := utils.AuthToken.GenerateAuthToken(jwt.SigningMethodRS256, jwt.MapClaims{
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

	jsonBytes, err:= json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error parsing request body in json bytes: %v", err)
	}
	
	request, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}
	
	request.Header.Set(HEADERS.Authorization, fmt.Sprintf("Bearer %s", authTokenString))
	request.Header.Set(HEADERS.ContentType, "application/json")
	request.Header.Set(HEADERS.Service, DISCORD_SERVICE)
	
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("error sending request to RDS Backend API: %v", err)
	}
	defer response.Body.Close()

	message:= "Something went wrong while generating verification link"
	verificationSiteURL := ""

	createMessage := func (format string) string {
		return fmt.Sprintf(format, VERIFICATION_STRING, verificationSiteURL, token, VERIFICATION_SUBSTRING)
	}

	if response.StatusCode == 201 || response.StatusCode == 200 {
		if metaData["dev"] == "true" {
			verificationSiteURL = config.AppConfig.MAIN_SITE_URL;
			message = createMessage("%s\n%s/discord?dev=true&token=%s\n%s")
		} else {
			verificationSiteURL = config.AppConfig.VERIFICATION_SITE_URL;
			message = createMessage("%s\n%s/discord?token=%s\n%s")
		}
	}

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
		return fmt.Errorf("error creating session: %v", err)
	}
<<<<<<< HEAD
=======
	
	request, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		return err
	}
	
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	request.Header.Set("Content-Type", "application/json")
	
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	
	defer response.Body.Close()

	message:= ""
	if response.StatusCode == 201 || response.StatusCode == 200 {
		verificationSiteURL := "";
		if CS.discordMessage.MetaData["dev"] == "true" {
			verificationSiteURL = config.AppConfig.MAIN_SITE_URL;
			message = fmt.Sprintf("%s\n%s/discord?dev=true&token=%s\n%s", VERIFICATION_STRING, verificationSiteURL, token, VERIFICATION_SUBSTRING)
		}
		if metaData["dev"] == "true" {
			verificationSiteURL = config.AppConfig.VERIFICATION_SITE_URL;
			message = fmt.Sprintf("%s\n%s/discord?dev=true&token=%s\n%s", VERIFICATION_STRING, verificationSiteURL, token, VERIFICATION_SUBSTRING)
		}
	}

	session, err := CreateSession()
	if err != nil {
		return err
	}
>>>>>>> 4f1bed6 (fix(verify): add logic to make a http call with website-backend and content generation)

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
