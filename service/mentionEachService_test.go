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
	// Save original SendMessage function and restore after tests
	originalSendMessage := queue.SendMessage
	defer func() {
		queue.SendMessage = originalSendMessage
	}()

	// Create role option - now a string instead of a map
	roleOption := &discordgo.ApplicationCommandInteractionDataOption{
		Name:  "role",
		Value: "123456789", // Now just the role ID as a string
	}

	// Create base Discord message
	mockData := &dtos.Data{
		GuildId: "876543210987654321",
		ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
			Name:    utils.CommandNames.MentionEach,
			Options: []*discordgo.ApplicationCommandInteractionDataOption{roleOption},
		},
	}

	discordMessage := &dtos.DiscordMessage{
		Data:      mockData,
		ChannelId: "987654321",
		Member: &discordgo.Member{
			User: &discordgo.User{
				ID: "user123",
			},
		},
	}

	t.Run("should queue message with role option only", func(t *testing.T) {
		// Mock queue.SendMessage
		var capturedPacket *dtos.DataPacket
		queue.SendMessage = func(message []byte) error {
			packetData := &dtos.DataPacket{}
			err := packetData.FromByte(message)
			assert.NoError(t, err)
			capturedPacket = packetData
			return nil
		}

		// Create request
		req, _ := http.NewRequest("POST", "/mention-each", bytes.NewBuffer([]byte("{}")))
		rr := httptest.NewRecorder()

		// Call service
		commandService := &CommandService{discordMessage: discordMessage}
		commandService.MentionEachService(rr, req)

		// Verify response - updated to match new message format
		assert.Contains(t, rr.Body.String(), "Mentioning all users with the 123456789")

		// Verify queued message - no role_name in metadata
		assert.NotNil(t, capturedPacket)
		assert.Equal(t, utils.CommandNames.MentionEach, capturedPacket.CommandName)
		assert.Equal(t, "user123", capturedPacket.UserID)
		assert.Equal(t, "123456789", capturedPacket.MetaData["role_id"])
		assert.Equal(t, "", capturedPacket.MetaData["message"])
		assert.Equal(t, "false", capturedPacket.MetaData["dev"])
		assert.Equal(t, "false", capturedPacket.MetaData["dev_title"])
	})

	t.Run("should include optional parameters when provided", func(t *testing.T) {
		// Add more options
		messageOption := &discordgo.ApplicationCommandInteractionDataOption{
			Name:  "message",
			Value: "Hello everyone!",
		}

		devOption := &discordgo.ApplicationCommandInteractionDataOption{
			Name:  "dev",
			Value: true,
		}

		// Update options
		mockData.Options = []*discordgo.ApplicationCommandInteractionDataOption{
			roleOption,
			messageOption,
			devOption,
		}

		// Mock queue.SendMessage
		var capturedPacket *dtos.DataPacket
		queue.SendMessage = func(message []byte) error {
			packetData := &dtos.DataPacket{}
			err := packetData.FromByte(message)
			assert.NoError(t, err)
			capturedPacket = packetData
			return nil
		}

		// Create request
		req, _ := http.NewRequest("POST", "/mention-each", bytes.NewBuffer([]byte("{}")))
		rr := httptest.NewRecorder()

		// Call service
		commandService := &CommandService{discordMessage: discordMessage}
		commandService.MentionEachService(rr, req)

		// Verify response - updated message
		assert.Contains(t, rr.Body.String(), "Sending individual mentions to users with the 123456789")

		// Verify queued message
		assert.NotNil(t, capturedPacket)
		assert.Equal(t, "Hello everyone!", capturedPacket.MetaData["message"])
		assert.Equal(t, "true", capturedPacket.MetaData["dev"])
	})

	t.Run("should handle queue errors", func(t *testing.T) {
		// Mock queue.SendMessage to fail
		queue.SendMessage = func(message []byte) error {
			return assert.AnError
		}

		// Create request
		req, _ := http.NewRequest("POST", "/mention-each", bytes.NewBuffer([]byte("{}")))
		rr := httptest.NewRecorder()

		// Call service
		commandService := &CommandService{discordMessage: discordMessage}
		commandService.MentionEachService(rr, req)

		// Verify error response - updated message
		assert.Contains(t, rr.Body.String(), "Failed to process your request")
	})

	t.Run("should handle missing role option", func(t *testing.T) {
		// Empty options
		mockData.Options = []*discordgo.ApplicationCommandInteractionDataOption{}

		// Create request
		req, _ := http.NewRequest("POST", "/mention-each", bytes.NewBuffer([]byte("{}")))
		rr := httptest.NewRecorder()

		// Call service
		commandService := &CommandService{discordMessage: discordMessage}
		commandService.MentionEachService(rr, req)

		assert.Contains(t, rr.Body.String(), "Role is required for mention-each command")
	})

	t.Run("should handle nil checks", func(t *testing.T) {

		brokenMessage := &dtos.DiscordMessage{
			Data: mockData,
		}

		req, _ := http.NewRequest("POST", "/mention-each", bytes.NewBuffer([]byte("{}")))
		rr := httptest.NewRecorder()

		commandService := &CommandService{discordMessage: brokenMessage}
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
