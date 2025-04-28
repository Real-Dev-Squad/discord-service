package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/queue"
	_ "github.com/Real-Dev-Squad/discord-service/tests/helpers"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func setupServiceTest(discordMessage *dtos.DiscordMessage) (*http.Request, *httptest.ResponseRecorder, *CommandService) {
	req, _ := http.NewRequest("POST", "/mention-each", bytes.NewBuffer([]byte("{}")))
	rr := httptest.NewRecorder()
	commandService := &CommandService{discordMessage: discordMessage}
	return req, rr, commandService
}

func TestMentionEachService(t *testing.T) {

	originalSendMessage := queue.SendMessage
	t.Cleanup(func() {
		queue.SendMessage = originalSendMessage
	})

	roleID := "123456789"

	createDefaultDiscordMessage := func(options []*discordgo.ApplicationCommandInteractionDataOption) *dtos.DiscordMessage {
		member := &discordgo.Member{
			User: &discordgo.User{ID: "user123"},
		}
		return &dtos.DiscordMessage{
			Data: &dtos.Data{
				GuildId: "guild123",
				ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
					Name:    utils.CommandNames.MentionEach,
					Options: options,
				},
			},
			ChannelId: "chan123",
			Member:    member,
		}
	}

	t.Run("should queue message when ff_enabled=true and only role provided", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{
			{Name: "role", Value: roleID},
			{Name: "ff_enabled", Value: true},
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
		req, rr, commandService := setupServiceTest(discordMessage)

		commandService.MentionEachService(rr, req)
		expectedSubString := "Mentioning all users with the \\u003c@\\u0026" + roleID + "\\u003e"
		assert.Contains(t, rr.Body.String(), expectedSubString)
		assert.NotNil(t, capturedPacket)
		assert.Equal(t, roleID, capturedPacket.MetaData["role_id"])
		assert.Equal(t, "false", capturedPacket.MetaData["dev"])
		assert.Equal(t, "false", capturedPacket.MetaData["dev_title"])
	})

	t.Run("should include optional params when ff_enabled=true", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{
			{Name: "role", Value: roleID},
			{Name: "message", Value: "Hello everyone!"},
			{Name: "dev", Value: true},
			{Name: "dev_title", Value: true},
			{Name: "ff_enabled", Value: true},
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
		req, rr, commandService := setupServiceTest(discordMessage)

		commandService.MentionEachService(rr, req)

		expectedSubString := "Fetching users with the \\u003c@\\u0026" + roleID + "\\u003e"
		assert.Contains(t, rr.Body.String(), expectedSubString)
		assert.NotNil(t, capturedPacket)
		assert.Equal(t, "Hello everyone!", capturedPacket.MetaData["message"])
		assert.Equal(t, "true", capturedPacket.MetaData["dev"])
		assert.Equal(t, "true", capturedPacket.MetaData["dev_title"])
	})

	t.Run("should return disabled error when ff_enabled option is missing", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{{Name: "role", Value: roleID}}
		discordMessage := createDefaultDiscordMessage(opts)

		queueCalled := false
		queue.SendMessage = func(message []byte) error { queueCalled = true; return nil }
		req, rr, commandService := setupServiceTest(discordMessage)

		commandService.MentionEachService(rr, req)

		assert.Contains(t, rr.Body.String(), "command requires the `ff_enabled:True` option")
		var resp discordgo.InteractionResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		if err == nil && resp.Data != nil {
			assert.Equal(t, discordgo.MessageFlagsEphemeral, resp.Data.Flags, "Error response should be ephemeral")
		}
		assert.False(t, queueCalled, "queue.SendMessage should NOT have been called")
	})

	t.Run("should return disabled error when ff_enabled option is false", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{
			{Name: "role", Value: roleID},
			{Name: "ff_enabled", Value: false},
		}
		discordMessage := createDefaultDiscordMessage(opts)

		queueCalled := false
		queue.SendMessage = func(message []byte) error { queueCalled = true; return nil }
		req, rr, commandService := setupServiceTest(discordMessage)

		commandService.MentionEachService(rr, req)

		assert.Contains(t, rr.Body.String(), "command requires the `ff_enabled:True` option")
		var resp discordgo.InteractionResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		if err == nil && resp.Data != nil {
			assert.Equal(t, discordgo.MessageFlagsEphemeral, resp.Data.Flags, "Error response should be ephemeral")
		}
		assert.False(t, queueCalled, "queue.SendMessage should NOT have been called")
	})

	t.Run("should handle queue errors when ff_enabled=true", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{
			{Name: "role", Value: roleID},
			{Name: "ff_enabled", Value: true},
		}
		discordMessage := createDefaultDiscordMessage(opts)
		queue.SendMessage = func(message []byte) error { return assert.AnError }
		req, rr, commandService := setupServiceTest(discordMessage)

		commandService.MentionEachService(rr, req)

		assert.Contains(t, rr.Body.String(), "Failed to process your request")
	})

	t.Run("should handle missing role option even if ff_enabled=true", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{
			{Name: "ff_enabled", Value: true},
		}
		discordMessage := createDefaultDiscordMessage(opts)
		req, rr, commandService := setupServiceTest(discordMessage)

		commandService.MentionEachService(rr, req)

		assert.Contains(t, rr.Body.String(), "Role is required")
	})

	t.Run("should handle invalid role format (non-string) when ff_enabled=true", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{
			{Name: "role", Value: 12345},
			{Name: "ff_enabled", Value: true},
		}
		discordMessage := createDefaultDiscordMessage(opts)
		req, rr, commandService := setupServiceTest(discordMessage)
		commandService.MentionEachService(rr, req)
		assert.Contains(t, rr.Body.String(), "Invalid role format")
	})

	t.Run("should handle invalid role format (empty string) when ff_enabled=true", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{
			{Name: "role", Value: ""},
			{Name: "ff_enabled", Value: true},
		}
		discordMessage := createDefaultDiscordMessage(opts)
		req, rr, commandService := setupServiceTest(discordMessage)
		commandService.MentionEachService(rr, req)
		assert.Contains(t, rr.Body.String(), "Invalid role format (empty ID)")
	})

	t.Run("should handle invalid ff_enabled type (defaults to disabled)", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{
			{Name: "role", Value: roleID},
			{Name: "ff_enabled", Value: "not-a-bool"},
		}
		discordMessage := createDefaultDiscordMessage(opts)

		queueCalled := false
		queue.SendMessage = func(message []byte) error { queueCalled = true; return nil }
		req, rr, commandService := setupServiceTest(discordMessage)

		commandService.MentionEachService(rr, req)

		assert.Contains(t, rr.Body.String(), "command requires the `ff_enabled:True` option")
		assert.False(t, queueCalled, "queue.SendMessage should NOT have been called")
	})

	t.Run("should handle nil checks when ff_enabled=true", func(t *testing.T) {
		opts := []*discordgo.ApplicationCommandInteractionDataOption{
			{Name: "role", Value: roleID},
			{Name: "ff_enabled", Value: true},
		}

		req, _ := http.NewRequest("POST", "/mention-each", bytes.NewBuffer([]byte("{}")))
		rr := httptest.NewRecorder()
		discordMessage := createDefaultDiscordMessage(opts)
		discordMessage.Data = nil
		commandService := &CommandService{discordMessage: discordMessage}
		commandService.MentionEachService(rr, req)
		assert.Contains(t, rr.Body.String(), "Invalid request data")

		rr = httptest.NewRecorder()
		discordMessage = createDefaultDiscordMessage(opts)
		discordMessage.Member = nil
		commandService = &CommandService{discordMessage: discordMessage}
		commandService.MentionEachService(rr, req)
		assert.Contains(t, rr.Body.String(), "Invalid request data")

		rr = httptest.NewRecorder()
		discordMessage = createDefaultDiscordMessage(opts)
		if discordMessage.Member == nil {
			discordMessage.Member = &discordgo.Member{}
		}
		discordMessage.Member.User = nil
		commandService = &CommandService{discordMessage: discordMessage}
		commandService.MentionEachService(rr, req)
		assert.Contains(t, rr.Body.String(), "Invalid request data")
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
