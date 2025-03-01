package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/fixtures"
	_ "github.com/Real-Dev-Squad/discord-service/tests/helpers"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestMainService(t *testing.T) {
	t.Run("should return Hello when command name is hello", func(t *testing.T) {
		CS := CommandService{discordMessage: fixtures.HelloCommand}
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		CS.HandleMessage(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should trigger ListeningService when command name is listening", func(t *testing.T) {
		discordMessage := dtos.DiscordMessage{
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

		CS := CommandService{discordMessage: discordMessage}
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		CS.HandleMessage(w, r)
		messageResponse := discordgo.InteractionResponse{}
		err := json.Unmarshal(w.Body.Bytes(), &messageResponse)
		assert.NoError(t, err)
		assert.Equal(t, messageResponse.Data.Content, "You are already set to listen.")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should fail if ListeningService returns an error", func(t *testing.T) {
		discordMessage := dtos.DiscordMessage{
			Member: &discordgo.Member{
				Nick: "test",
				User: &discordgo.User{
					ID: "123456789012345678",
				},
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
		config.AppConfig.MAX_RETRIES = 1
		CS := CommandService{discordMessage: discordMessage}
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		CS.HandleMessage(w, r)
		response := dtos.CommandError{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Message, "Queue channel is not initialized")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("should trigger VerifyService when command name is verify", func(t *testing.T) {
		discordMessage := dtos.DiscordMessage{
			Data: &dtos.Data{
				GuildId: "876543210987654321",
				ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
					Name: utils.CommandNames.Verify,
				},
			},
		}

		CS := CommandService{discordMessage: discordMessage}
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		CS.HandleMessage(w, r)
		messageResponse := discordgo.InteractionResponse{}
		err := json.Unmarshal(w.Body.Bytes(), &messageResponse)
		assert.NoError(t, err)
		assert.Equal(t, messageResponse.Data.Content, "Work in progress for Verify command")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return default handler when command name is not in record", func(t *testing.T) {
		discordMessage := dtos.DiscordMessage{
			Data: &dtos.Data{
				GuildId: "876543210987654321",
				ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
					Name: "unknown",
				},
			},
		}

		CS := CommandService{discordMessage: discordMessage}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		CS.HandleMessage(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}
