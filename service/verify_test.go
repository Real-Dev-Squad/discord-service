package service

import (
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

	t.Run("should return success response with 200 status code", func(t *testing.T) {
		message := *baseMessage
		message.Data = &dtos.Data{}
		service := &CommandService{
			discordMessage: &message,
		}

		queue.SendMessage = func(data []byte) error { return nil }

		r, _ := http.NewRequest("POST", "/verify", nil)
		w := httptest.NewRecorder()

		service.Verify(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return internal server error when queue send message fails", func(t *testing.T) {
		message := *baseMessage
		message.Data = &dtos.Data{}
		service := &CommandService{
			discordMessage: &message,
		}
		queue.SendMessage = func(data []byte) error {
			return errors.New("queue error")
		}

		r, _ := http.NewRequest("POST", "/verify", nil)
		w := httptest.NewRecorder()
		service.Verify(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return internal server error when json marshalling fails", func(t *testing.T) {
		message := *baseMessage
		message.Data = &dtos.Data{}

		utils.Json = &mockFailingJsonHandler{}
		service := &CommandService{
			discordMessage: &message,
		}

		r, _ := http.NewRequest("POST", "/verify", nil)
		w := httptest.NewRecorder()
		service.Verify(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
