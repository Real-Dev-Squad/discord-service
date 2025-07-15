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
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/golang-jwt/jwt/v5"
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

func (m *mockDiscordSession) Close() error {
	return nil
}

type mockFailingDiscordSession struct {
	*discordgo.Session
}

func (m *mockFailingDiscordSession) WebhookMessageEdit(webhookID, token, messageID string, data *discordgo.WebhookEdit, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	return nil, errors.New("webhook error")
}

type mockFailingDiscordSessionCloser struct {
	*discordgo.Session
}

func (m *mockFailingDiscordSessionCloser) WebhookMessageEdit(webhookID, token, messageID string, data *discordgo.WebhookEdit, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	return nil, nil
}

func (m *mockFailingDiscordSessionCloser) Close() error {
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

type mockFailingUniqueToken struct{}

func (m *mockFailingUniqueToken) GenerateUniqueToken() (string, error) {
	return "", errors.New("unique token error")
}

type mockFailingAuthToken struct{}

func (m *mockFailingAuthToken) GenerateAuthToken(method jwt.SigningMethod, claims jwt.Claims, privateKey any) (string, error) {
	return "", errors.New("auth token error")
}

func TestVerify(t *testing.T) {
	privateKey := generateTestPrivateKey(t)
	pemPrivateKey := pemEncodePrivateKey(privateKey)

	originalCreateSession := CreateSession
	originalBotPrivateKey := config.AppConfig.BOT_PRIVATE_KEY
	originalUniqueToken := utils.UniqueToken
	originalAuthToken := utils.AuthToken
	config.AppConfig.BOT_PRIVATE_KEY = pemPrivateKey

	t.Run("should return error when fails to generate unique token", func(t *testing.T) {
		handler := &CommandHandler{discordMessage: &dtos.DataPacket{}}
		utils.UniqueToken = &mockFailingUniqueToken{}
		t.Cleanup(func() {
			utils.UniqueToken = originalUniqueToken
		})
		err := handler.verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error generating unique token")
	})

	t.Run("should return error when fails to parse private key string to rsa private key", func(t *testing.T) {
		config.AppConfig.BOT_PRIVATE_KEY = "<invalid-key>"
		t.Cleanup(func() {
			config.AppConfig.BOT_PRIVATE_KEY = originalBotPrivateKey
		})
		handler := &CommandHandler{
			discordMessage: &dtos.DataPacket{},
		}
		err := handler.verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error parsing private key string to rsa private key")
	})

	t.Run("should return error when fails to generate auth token", func(t *testing.T) {
		config.AppConfig.BOT_PRIVATE_KEY = pemPrivateKey
		handler := &CommandHandler{
			discordMessage: &dtos.DataPacket{},
		}
		utils.AuthToken = &mockFailingAuthToken{}
		t.Cleanup(func() {
			utils.AuthToken = originalAuthToken
		})
		err := handler.verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error generating auth token")
	})

	t.Run("should return error when fails to create http request", func(t *testing.T) {
		config.AppConfig.BOT_PRIVATE_KEY = pemPrivateKey
		config.AppConfig.RDS_BASE_API_URL = "http://localhost:1234\x7f"
		handler := &CommandHandler{
			discordMessage: &dtos.DataPacket{},
		}
		err := handler.verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error creating http request")
	})

	t.Run("should return error when fails to send request to RDS Backend API", func(t *testing.T) {
		config.AppConfig.BOT_PRIVATE_KEY = pemPrivateKey
		config.AppConfig.RDS_BASE_API_URL = "http://localhost:12345"
		handler := &CommandHandler{
			discordMessage: &dtos.DataPacket{},
		}
		err := handler.verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error sending request to RDS Backend API")
	})

	t.Run("should return error when fails to create session", func(t *testing.T) {
		config.AppConfig.BOT_PRIVATE_KEY = pemPrivateKey
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		config.AppConfig.RDS_BASE_API_URL = server.URL
		CreateSession = func() (DiscordSessionWrapper, error) {
			return nil, errors.New("session error")
		}
		t.Cleanup(func() {
			CreateSession = originalCreateSession
		})

		handler := &CommandHandler{
			discordMessage: &dtos.DataPacket{},
		}

		err := handler.verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error creating session")
	})

	t.Run("should return error when fails to edit original message", func(t *testing.T) {
		config.AppConfig.BOT_PRIVATE_KEY = pemPrivateKey
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		config.AppConfig.RDS_BASE_API_URL = server.URL

		CreateSession = func() (DiscordSessionWrapper, error) {
			return &mockFailingDiscordSession{}, nil
		}
		t.Cleanup(func() {
			CreateSession = originalCreateSession
		})

		handler := &CommandHandler{discordMessage: &dtos.DataPacket{}}
		err := handler.verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error editing original message for application")
	})

	t.Run("should return error when fails to close session", func(t *testing.T) {
		config.AppConfig.BOT_PRIVATE_KEY = pemPrivateKey
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		config.AppConfig.RDS_BASE_API_URL = server.URL

		CreateSession = func() (DiscordSessionWrapper, error) {
			return &mockFailingDiscordSessionCloser{}, nil
		}
		t.Cleanup(func() {
			CreateSession = originalCreateSession
		})
		handler := &CommandHandler{discordMessage: &dtos.DataPacket{}}
		err := handler.verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error closing session")
	})

	t.Run("should not return error when succeeds", func(t *testing.T) {
		config.AppConfig.BOT_PRIVATE_KEY = pemPrivateKey
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		config.AppConfig.RDS_BASE_API_URL = server.URL

		CreateSession = func() (DiscordSessionWrapper, error) {
			return &mockDiscordSession{}, nil
		}
		t.Cleanup(func() {
			CreateSession = originalCreateSession
		})

		handler := &CommandHandler{discordMessage: &dtos.DataPacket{
			MetaData: map[string]string{
				"dev": "true",
			},
		}}
		
		err := handler.verify()
		assert.NoError(t, err)
	})
}
