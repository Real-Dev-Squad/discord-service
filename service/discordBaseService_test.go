package service

import (
	"bytes"
	"encoding/json"
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
		rr := httptest.NewRecorder()
		DiscordBaseService(rr, r)
		bytes, err := json.Marshal(map[string]string{
			"error": "Invalid Request Payload",
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, string(bytes) + "\n", rr.Body.String())
	})

	t.Run("should return internal server error when fails to unmarshal request payload", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(`"name": "listening"`)))
		rr := httptest.NewRecorder()
		DiscordBaseService(rr, r)
		bytes, err := json.Marshal(map[string]string{
			"error": "Internal Server Error",
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, string(bytes) + "\n", rr.Body.String())
	})

	t.Run("should return 400 status code when message type is unknown", func(t *testing.T) {
		msgByte, err := json.Marshal(dtos.DiscordMessage{Type: 3})
		assert.NoError(t, err)
		r := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(msgByte)))
		rr := httptest.NewRecorder()
		DiscordBaseService(rr, r)
		bytes, err := json.Marshal(map[string]string{
			"error": "Invalid Request Payload",
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, string(bytes) + "\n", rr.Body.String())
	})

	t.Run("should return success response when message type is interaction ping", func(t *testing.T) {
		msgByte, err := json.Marshal(map[string]interface{}{"type": discordgo.InteractionPing, "message": "Pong"})
		assert.NoError(t, err)
		r := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(msgByte)))
		rr := httptest.NewRecorder()
		DiscordBaseService(rr, r)
		bytes, err:= json.Marshal(map[string]any{
			"data": discordgo.InteractionResponsePong,
			"message": "Pong",
		})
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, string(bytes) + "\n", rr.Body.String())
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
