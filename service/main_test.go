package service

import (
	"encoding/json"
	"fmt"
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

	t.Run("should return ListeningService when command name is listening", func(t *testing.T) {
		discordMessage := &dtos.DiscordMessage{
			Member: &discordgo.Member{
				Nick: fmt.Sprintf("test%s", utils.NICKNAME_SUFFIX),
			},
			Data: &dtos.Data{
				GuildId: "876543210987654321",
				ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
					Name: "listening",
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
