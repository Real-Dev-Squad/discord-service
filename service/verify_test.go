package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/queue"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
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
		originalSendMessage := queue.SendMessage
		defer func() {
			queue.SendMessage = originalSendMessage
		}()

		message := *baseMessage
		message.Data = &dtos.Data{
			ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "dev",
						Value: true,
					},
				},
			},
		}
		service := &CommandService{
			discordMessage: &message,
		}

		queue.SendMessage = func(data []byte) error { return nil }

		req := httptest.NewRequest("POST", "/verify", nil)
		rr := httptest.NewRecorder()

		service.Verify(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return internal server error when queue send message fails", func(t *testing.T) {
		originalSendMessage := queue.SendMessage
		defer func() {
			queue.SendMessage = originalSendMessage
		}()

		message := *baseMessage
		message.Data = &dtos.Data{}
		service := &CommandService{
			discordMessage: &message,
		}
		queue.SendMessage = func(data []byte) error {
			return errors.New("queue error")
		}

		req := httptest.NewRequest("POST", "/verify", nil)
		rr := httptest.NewRecorder()
		service.Verify(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}