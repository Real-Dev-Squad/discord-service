package handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

type mockDiscordSession struct {
	*discordgo.Session
	capturedMessage *string
}

func (m *mockDiscordSession) WebhookMessageEdit(webhookID, token, messageID string, data *discordgo.WebhookEdit, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	if m.capturedMessage != nil {
		*m.capturedMessage = *data.Content
	}
	return nil, nil
}

func (m *mockDiscordSession) GuildMemberNickname(guildID, userID, nickname string, options ...discordgo.RequestOption) error {
	return nil
}

func (m *mockDiscordSession) Close() error {
	return nil
}

type mockFailingDiscordSession struct {
	*discordgo.Session
}

func (m *mockFailingDiscordSession) WebhookMessageEdit(webhookID, token, messageID string, data *discordgo.WebhookEdit, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	return nil, errors.New("webhook error")
}

func (m *mockFailingDiscordSession) GuildMemberNickname(guildID, userID, nickname string, options ...discordgo.RequestOption) error {
	return errors.New("nickname error")
}

func (m *mockFailingDiscordSession) Close() error {
	return errors.New("close error")
}

func generateTestPrivateKey(t *testing.T) *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	return privateKey
}

func pemEncodePrivateKey(privateKey *rsa.PrivateKey) string {
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}))
}

func TestVerify(t *testing.T) {
	privateKey := generateTestPrivateKey(t)
	pemPrivateKey := pemEncodePrivateKey(privateKey)

	originalCreateSession := CreateSession
	originalBotPrivateKey := config.AppConfig.BOT_PRIVATE_KEY
	defer func() {
		CreateSession = originalCreateSession
		config.AppConfig.BOT_PRIVATE_KEY = originalBotPrivateKey
	}()

	config.AppConfig.BOT_PRIVATE_KEY = pemPrivateKey

	t.Run("success in dev mode", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}))
		defer server.Close()
		config.AppConfig.RDS_BASE_API_URL = server.URL
		config.AppConfig.VERIFICATION_SITE_URL = "http://dev.realdevsquad.com"

		var capturedMessage string
		CreateSession = func() (DiscordSessionWrapper, error) {
			return &mockDiscordSession{Session: &discordgo.Session{}, capturedMessage: &capturedMessage}, nil
		}

		handler := &CommandHandler{
			discordMessage: &dtos.DataPacket{
				UserID: "userID",
				MetaData: map[string]string{
					"dev":             "true",
					"applicationId":   "appID",
					"token":           "discord-token",
					"userAvatarHash":  "avatarHash",
					"userName":        "testuser",
					"discriminator":   "1234",
					"discordJoinedAt": "somedate",
				},
			},
		}

		err := handler.verify()
		assert.NoError(t, err)
		assert.Contains(t, capturedMessage, VERIFICATION_STRING)
	})

	t.Run("error on http request", func(t *testing.T) {
		config.AppConfig.RDS_BASE_API_URL = "http://localhost:12345"
		handler := &CommandHandler{
			discordMessage: &dtos.DataPacket{
				UserID: "userID",
				MetaData: map[string]string{
					"dev":             "true",
					"applicationId":   "appID",
					"token":           "discord-token",
					"userAvatarHash":  "avatarHash",
					"userName":        "testuser",
					"discriminator":   "1234",
					"discordJoinedAt": "somedate",
				},
			},
		}
		err := handler.verify()
		assert.Error(t, err)
	})

	t.Run("error on create session", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}))
		defer server.Close()
		config.AppConfig.RDS_BASE_API_URL = server.URL

		CreateSession = func() (DiscordSessionWrapper, error) {
			return nil, errors.New("session error")
		}

		handler := &CommandHandler{
			discordMessage: &dtos.DataPacket{
				UserID: "userID",
				MetaData: map[string]string{
					"dev":             "true",
					"applicationId":   "appID",
					"token":           "discord-token",
					"userAvatarHash":  "avatarHash",
					"userName":        "testuser",
					"discriminator":   "1234",
					"discordJoinedAt": "somedate",
				},
			},
		}
		err := handler.verify()
		assert.Error(t, err)
		assert.Equal(t, "session error", err.Error())
	})

	t.Run("error on webhook message edit", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		config.AppConfig.RDS_BASE_API_URL = server.URL
		config.AppConfig.MAIN_SITE_URL = "http://realdevsquad.com"

		CreateSession = func() (DiscordSessionWrapper, error) {
			// Initialize with an empty but non-nil session to avoid panics
			return &mockFailingDiscordSession{Session: &discordgo.Session{}}, nil
		}

		handler := &CommandHandler{
			discordMessage: &dtos.DataPacket{
				UserID: "userID",
				MetaData: map[string]string{
					"dev":             "true",
					"applicationId":   "appID",
					"token":           "discord-token",
					"userAvatarHash":  "avatarHash",
					"userName":        "testuser",
					"discriminator":   "1234",
					"discordJoinedAt": "somedate",
				},
			},
		}
		err := handler.verify()
		assert.Error(t, err)
		assert.Equal(t, "webhook error", err.Error())
	})
}
