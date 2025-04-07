package service

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/queue"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestMentionEachService(t *testing.T) {

	originalSendMessage := queue.SendMessage
	defer func() {
		queue.SendMessage = originalSendMessage
	}()

	roleID := "123456789"
	createDefaultDiscordMessage := func(options []*discordgo.ApplicationCommandInteractionDataOption) *dtos.DiscordMessage {
		return &dtos.DiscordMessage{
			Data: &dtos.Data{
				GuildId: "876543210987654321",
				ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
					Name:    utils.CommandNames.MentionEach,
					Options: options,
				},
			},
			ChannelId: "987654321",
			Member: &discordgo.Member{
				User: &discordgo.User{ID: "user123"},
			},
		}
	}

	t.Run("should queue message with role option only", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{
			{Name: "role", Value: roleID},
		}
		discordMessage := createDefaultDiscordMessage(opts)
		var capturedPacket *dtos.DataPacket
		queue.SendMessage = func(message []byte) error {
			packetData := &dtos.DataPacket{}
			err := packetData.FromByte(message)
			assert.NoError(t, err)
			capturedPacket = packetData
			return nil
		}

		req, _ := http.NewRequest("POST", "/mention-each", bytes.NewBuffer([]byte("{}")))
		rr := httptest.NewRecorder()
		commandService := &CommandService{discordMessage: discordMessage}
		commandService.MentionEachService(rr, req)
		expectedSubString := "Mentioning all users with the \\u003c@\\u0026" + roleID + "\\u003e"
		assert.Contains(t, rr.Body.String(), expectedSubString)
		assert.NotNil(t, capturedPacket)
		assert.Equal(t, utils.CommandNames.MentionEach, capturedPacket.CommandName)
		assert.Equal(t, "user123", capturedPacket.UserID)
		assert.Equal(t, "123456789", capturedPacket.MetaData["role_id"])
		assert.Equal(t, "", capturedPacket.MetaData["message"])
		assert.Equal(t, "false", capturedPacket.MetaData["dev"])
		assert.Equal(t, "false", capturedPacket.MetaData["dev_title"])
	})

	t.Run("should include optional parameters when provided", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{
			{Name: "role", Value: roleID},
			{Name: "message", Value: "Hello everyone!"},
			{Name: "dev", Value: true},
			{Name: "dev_title", Value: true},
		}
		discordMessage := createDefaultDiscordMessage(opts)

		var capturedPacket *dtos.DataPacket
		queue.SendMessage = func(message []byte) error {
			packetData := &dtos.DataPacket{}
			err := packetData.FromByte(message)
			assert.NoError(t, err)
			capturedPacket = packetData
			return nil
		}

		req, _ := http.NewRequest("POST", "/mention-each", bytes.NewBuffer([]byte("{}")))
		rr := httptest.NewRecorder()

		commandService := &CommandService{discordMessage: discordMessage}
		commandService.MentionEachService(rr, req)
		expectedSubString := "Fetching users with the \\u003c@\\u0026" + roleID + "\\u003e"

		assert.Contains(t, rr.Body.String(), expectedSubString)
		assert.NotNil(t, capturedPacket)
		assert.Equal(t, "Hello everyone!", capturedPacket.MetaData["message"])
		assert.Equal(t, "true", capturedPacket.MetaData["dev"])
	})

	t.Run("should handle queue errors", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{{Name: "role", Value: roleID}}
		discordMessage := createDefaultDiscordMessage(opts)
		queue.SendMessage = func(message []byte) error {
			return assert.AnError
		}

		req, _ := http.NewRequest("POST", "/mention-each", bytes.NewBuffer([]byte("{}")))
		rr := httptest.NewRecorder()

		commandService := &CommandService{discordMessage: discordMessage}
		commandService.MentionEachService(rr, req)
		assert.Contains(t, rr.Body.String(), "Failed to process your request")
	})

	t.Run("should handle missing role option", func(t *testing.T) {
		discordMessage := createDefaultDiscordMessage([]*discordgo.ApplicationCommandInteractionDataOption{})

		req, _ := http.NewRequest("POST", "/mention-each", bytes.NewBuffer([]byte("{}")))
		rr := httptest.NewRecorder()

		commandService := &CommandService{discordMessage: discordMessage}
		commandService.MentionEachService(rr, req)

		assert.Contains(t, rr.Body.String(), "Role is required for mention-each command")
	})

	t.Run("should handle nil checks", func(t *testing.T) {
		// Test with nil discordMessage
		req, _ := http.NewRequest("POST", "/mention-each", bytes.NewBuffer([]byte("{}")))
		rr := httptest.NewRecorder()
		commandService := &CommandService{discordMessage: nil}
		commandService.MentionEachService(rr, req)
		assert.Contains(t, rr.Body.String(), "Invalid request data")

		// Test with nil Data
		rr = httptest.NewRecorder()
		discordMessage := createDefaultDiscordMessage(nil)
		discordMessage.Data = nil
		commandService = &CommandService{discordMessage: discordMessage}
		commandService.MentionEachService(rr, req)
		assert.Contains(t, rr.Body.String(), "Invalid request data")

		// Test with nil Member
		rr = httptest.NewRecorder()
		discordMessage = createDefaultDiscordMessage(nil)
		discordMessage.Member = nil
		commandService = &CommandService{discordMessage: discordMessage}
		commandService.MentionEachService(rr, req)
		assert.Contains(t, rr.Body.String(), "Invalid request data")

		// Test with nil User
		rr = httptest.NewRecorder()
		discordMessage = createDefaultDiscordMessage(nil)
		discordMessage.Member.User = nil
		commandService = &CommandService{discordMessage: discordMessage}
		commandService.MentionEachService(rr, req)
		assert.Contains(t, rr.Body.String(), "Invalid request data")
	})

	t.Run("should handle invalid role format (non-string)", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{
			{Name: "role", Value: 12345},
		}
		discordMessage := createDefaultDiscordMessage(opts)

		req, _ := http.NewRequest("POST", "/mention-each", bytes.NewBuffer([]byte("{}")))
		rr := httptest.NewRecorder()
		commandService := &CommandService{discordMessage: discordMessage}

		commandService.MentionEachService(rr, req)
		assert.Contains(t, rr.Body.String(), "Invalid role format")
	})

	t.Run("should set correct response content for dev_title=true", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{
			{Name: "role", Value: roleID},
			{Name: "dev_title", Value: true},
		}
		discordMessage := createDefaultDiscordMessage(opts)

		queue.SendMessage = func(message []byte) error { return nil }

		req, _ := http.NewRequest("POST", "/mention-each", bytes.NewBuffer([]byte("{}")))
		rr := httptest.NewRecorder()
		commandService := &CommandService{discordMessage: discordMessage}

		commandService.MentionEachService(rr, req)
		expectedSubString := "Fetching users with the \\u003c@\\u0026" + roleID + "\\u003e"
		assert.Contains(t, rr.Body.String(), expectedSubString)
	})

	t.Run("should set correct response content for dev=true", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{
			{Name: "role", Value: roleID},
			{Name: "message", Value: "Dev message"},
			{Name: "dev", Value: true},
		}
		discordMessage := createDefaultDiscordMessage(opts)
		var capturedPacket *dtos.DataPacket
		queue.SendMessage = func(message []byte) error {
			packetData := &dtos.DataPacket{}
			err := packetData.FromByte(message)
			assert.NoError(t, err)
			capturedPacket = packetData
			return nil
		}

		req, _ := http.NewRequest("POST", "/mention-each", bytes.NewBuffer([]byte("{}")))
		rr := httptest.NewRecorder()
		commandService := &CommandService{discordMessage: discordMessage}
		commandService.MentionEachService(rr, req)

		expectedSubString := "Sending individual mentions to users with the \\u003c@\\u0026" + roleID + "\\u003e"
		assert.Contains(t, rr.Body.String(), expectedSubString)
		assert.NotNil(t, capturedPacket)
		assert.Equal(t, "Dev message", capturedPacket.MetaData["message"])
		assert.Equal(t, "true", capturedPacket.MetaData["dev"])
		assert.Equal(t, "false", capturedPacket.MetaData["dev_title"])

	})

}

func TestFindOption(t *testing.T) {
	options := []*discordgo.ApplicationCommandInteractionDataOption{
		{Name: "role", Value: "role-id-123"},
		{Name: "message", Value: "Hello!"},
	}

	t.Run("should find option when present", func(t *testing.T) {
		option := findOption(options, "message")
		assert.NotNil(t, option)
		assert.Equal(t, "message", option.Name)
		assert.Equal(t, "Hello!", option.Value)
	})

	t.Run("should return nil when option not found", func(t *testing.T) {
		option := findOption(options, "nonexistent")
		assert.Nil(t, option)
	})

	t.Run("should handle empty options", func(t *testing.T) {
		option := findOption([]*discordgo.ApplicationCommandInteractionDataOption{}, "role")
		assert.Nil(t, option)
	})
}
