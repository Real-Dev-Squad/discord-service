package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/queue"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

type mockFailingJsonHandler struct{}

func (m *mockFailingJsonHandler) ToJson(data interface{}) (string, error) {
	return "", errors.New("json marshal error")
}

func TestVerify(t *testing.T) {
	originalSendMessage := queue.SendMessage
	originalJson := utils.Json
	defer func() {
		queue.SendMessage = originalSendMessage
		utils.Json = originalJson
	}()

	joinedAt, _ := time.Parse(time.RFC3339, "2022-01-01T00:00:00Z")

	baseMessage := &dtos.DiscordMessage{
		Member: &discordgo.Member{
			User: &discordgo.User{
				ID:            "userID-123",
				Username:      "testuser",
				Discriminator: "1234",
			},
			Avatar:   "avatar-hash",
			JoinedAt: joinedAt,
		},
		ChannelId:     "channelID-456",
		Token:         "interaction-token",
		ApplicationId: "appID-789",
	}

	t.Run("Success with dev as true", func(t *testing.T) {
		message := *baseMessage
		message.Data = &dtos.Data{
			ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
				Name: "verify",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "dev",
						Type:  discordgo.ApplicationCommandOptionBoolean,
						Value: true,
					},
				},
			},
		}

		service := &CommandService{
			discordMessage: &message,
		}

		var capturedData []byte
		queue.SendMessage = func(data []byte) error {
			capturedData = data
			return nil
		}

		req, _ := http.NewRequest("POST", "/verify", nil)
		rr := httptest.NewRecorder()
		service.Verify(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var packet dtos.DataPacket
		json.Unmarshal(capturedData, &packet)

		assert.Equal(t, "true", packet.MetaData["dev"])
		assert.Equal(t, "userID-123", packet.UserID)
	})

	t.Run("Success with dev as false", func(t *testing.T) {
		message := *baseMessage
		message.Data = &dtos.Data{
			ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
				Name: "verify",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "dev",
						Type:  discordgo.ApplicationCommandOptionBoolean,
						Value: false,
					},
				},
			},
		}
		service := &CommandService{
			discordMessage: &message,
		}
		var capturedData []byte
		queue.SendMessage = func(data []byte) error {
			capturedData = data
			return nil
		}

		req, _ := http.NewRequest("POST", "/verify", nil)
		rr := httptest.NewRecorder()
		service.Verify(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var packet dtos.DataPacket
		json.Unmarshal(capturedData, &packet)
		assert.Equal(t, "false", packet.MetaData["dev"])
	})

	t.Run("Success with no dev flag", func(t *testing.T) {
		message := *baseMessage
		message.Data = &dtos.Data{
			ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
				Name:    "verify",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{},
			},
		}
		service := &CommandService{
			discordMessage: &message,
		}

		var capturedData []byte
		queue.SendMessage = func(data []byte) error {
			capturedData = data
			return nil
		}

		req, _ := http.NewRequest("POST", "/verify", nil)
		rr := httptest.NewRecorder()
		service.Verify(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var packet dtos.DataPacket
		json.Unmarshal(capturedData, &packet)
		assert.Equal(t, "false", packet.MetaData["dev"])
	})

	t.Run("Queue send message fails", func(t *testing.T) {
		message := *baseMessage
		message.Data = &dtos.Data{
			ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
				Name: "verify",
			},
		}
		service := &CommandService{
			discordMessage: &message,
		}

		queue.SendMessage = func(data []byte) error {
			return errors.New("queue error")
		}

		req, _ := http.NewRequest("POST", "/verify", nil)
		rr := httptest.NewRecorder()
		service.Verify(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var resp dtos.Response
		json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Equal(t, "Something went wrong", resp.Message)
	})

	t.Run("Json marshalling fails", func(t *testing.T) {
		message := *baseMessage
		message.Data = &dtos.Data{
			ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
				Name: "verify",
			},
		}
		service := &CommandService{
			discordMessage: &message,
		}
		queue.SendMessage = func(data []byte) error {
			return nil
		}
		utils.Json = &mockFailingJsonHandler{}

		req, _ := http.NewRequest("POST", "/verify", nil)
		rr := httptest.NewRecorder()
		service.Verify(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, "Internal Server Error\n", rr.Body.String())
	})
}
