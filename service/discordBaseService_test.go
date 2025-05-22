package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestDiscordBaseService(t *testing.T) {
	t.Run("should return bad request error for empty request body", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", http.NoBody)
		w := httptest.NewRecorder()
		DiscordBaseService(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"success": false, "message": "Invalid Request Payload", "status": 400}` + "\n", w.Body.String())
	})

	t.Run("should return bad request error when fails to unmarshal request payload", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(`"name": "listening"`)))
		w := httptest.NewRecorder()
		DiscordBaseService(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"success": false, "message": "Invalid Request Payload", "status": 400}`+ "\n", w.Body.String())
	})

	t.Run("should return 400 status code when message type is unknown", func(t *testing.T) {
		msgByte, err := json.Marshal(dtos.DiscordMessage{Type: 3})
		assert.NoError(t, err)
		r := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(msgByte)))
		w := httptest.NewRecorder()
		DiscordBaseService(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return success response when message type is interaction ping", func(t *testing.T) {
		msgByte, err := json.Marshal(dtos.DiscordMessage{Type: discordgo.InteractionPing})
		assert.NoError(t, err)
		r := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(msgByte)))
		w := httptest.NewRecorder()
		DiscordBaseService(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fmt.Sprintf(`{"type":%d}`+"\n", discordgo.InteractionResponsePong), w.Body.String())
	})

	t.Run("should return success response when message type is interaction application command", func(t *testing.T) {
		msgByte, err := json.Marshal(dtos.DiscordMessage{
			Type: discordgo.InteractionApplicationCommand,
			Data: &dtos.Data{
				GuildId: "1234",
			},
		})
		assert.NoError(t, err)
		r := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(msgByte)))
		w := httptest.NewRecorder()
		DiscordBaseService(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
