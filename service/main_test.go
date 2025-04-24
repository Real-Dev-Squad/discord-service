package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Real-Dev-Squad/discord-service/queue"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/fixtures"
	_ "github.com/Real-Dev-Squad/discord-service/tests/helpers"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestMainService(t *testing.T) {
	t.Run("should return HelloService when command name is hello", func(t *testing.T) {
		handler := MainService(fixtures.HelloCommand)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		handler(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should trigger ListeningService when command name is listening", func(t *testing.T) {
		discordMessage := &dtos.DiscordMessage{
			Member: &discordgo.Member{
				Nick: fmt.Sprintf("test%s", utils.NICKNAME_SUFFIX),
			},
			Data: &dtos.Data{
				GuildId: "876543210987654321",
				ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
					Name: utils.CommandNames.Listening,
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{Value: true},
					},
				},
			},
		}

		handler := MainService(discordMessage)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		handler(w, r)
		messageResponse := discordgo.InteractionResponse{}
		err := json.Unmarshal(w.Body.Bytes(), &messageResponse)
		assert.NoError(t, err)
		assert.Equal(t, messageResponse.Data.Content, "You are already set to listen.")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should trigger VerifyService when command name is verify", func(t *testing.T) {
		discordMessage := &dtos.DiscordMessage{
			Data: &dtos.Data{
				GuildId: "876543210987654321",
				ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
					Name: utils.CommandNames.Verify,
				},
			},
		}

		handler := MainService(discordMessage)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		handler(w, r)
		messageResponse := discordgo.InteractionResponse{}
		err := json.Unmarshal(w.Body.Bytes(), &messageResponse)
		assert.NoError(t, err)
		assert.Equal(t, messageResponse.Data.Content, "Work in progress for Verify command")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should trigger MentionEachService when command name is mention-each", func(t *testing.T) {
		roleID := "role987"
		discordMessage := &dtos.DiscordMessage{
			Data: &dtos.Data{
				GuildId: "guild123",
				ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
					Name: utils.CommandNames.MentionEach,
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{Name: "role", Value: roleID},
					},
				},
			},
			Member:    &discordgo.Member{User: &discordgo.User{ID: "mentionUser"}, Permissions: discordgo.PermissionMentionEveryone}, // Need member + perms
			ChannelId: "chan123",
		}

		queueCalled := false
		queue.SendMessage = func(message []byte) error {
			queueCalled = true
			return nil
		}

		handler := MainService(discordMessage)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/", bytes.NewBuffer([]byte{}))
		handler(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		expectedSubString := "Mentioning all users with the \\u003c@\\u0026" + roleID + "\\u003e"
		assert.Contains(t, w.Body.String(), expectedSubString)
		assert.True(t, queueCalled, "queue.SendMessage should have been called for mention-each")
	})

	t.Run("should return default handler when command name is not in record", func(t *testing.T) {
		discordMessage := &dtos.DiscordMessage{
			Data: &dtos.Data{
				GuildId: "876543210987654321",
				ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
					Name: "unknown",
				},
			},
		}

		handler := MainService(discordMessage)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		handler(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

}
