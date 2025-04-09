package utils

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDiscordSession struct {
	mock.Mock
}

func (m *MockDiscordSession) GuildMembers(guildID, after string, limit int) ([]*discordgo.Member, error) {
	args := m.Called(guildID, after, limit)
	return args.Get(0).([]*discordgo.Member), args.Error(1)
}

func (m *MockDiscordSession) ChannelMessageSend(channelID, content string) (*discordgo.Message, error) {
	args := m.Called(channelID, content)
	return args.Get(0).(*discordgo.Message), args.Error(1)
}

func TestGetUsersWithRole(t *testing.T) {
	guildID := "testGuild"
	roleID := "testRole"

	t.Run("returns users with matching role", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		members := []*discordgo.Member{
			{User: &discordgo.User{ID: "123"}, Roles: []string{"testRole"}},
			{User: &discordgo.User{ID: "456"}, Roles: []string{"otherRole"}},
		}
		mockSession.On("GuildMembers", guildID, "", 1000).Return(members, nil)
		mockSession.On("ChannelMessageSend", mock.Anything, mock.Anything).Return(&discordgo.Message{}, nil)

		result, err := GetUsersWithRole(mockSession, guildID, roleID)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "123", result[0].User.ID)
	})

	t.Run("handles error from GuildMembers", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		mockSession.On("GuildMembers", guildID, "", 1000).Return([]*discordgo.Member{}, errors.New("API error"))
		mockSession.On("ChannelMessageSend", mock.Anything, mock.Anything).Return(&discordgo.Message{}, nil)

		_, err := GetUsersWithRole(mockSession, guildID, roleID)
		assert.Error(t, err)
	})

	t.Run("returns empty slice when no users have the role", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		members := []*discordgo.Member{
			{User: &discordgo.User{ID: "123"}, Roles: []string{"otherRole"}},
		}
		mockSession.On("GuildMembers", guildID, "", 1000).Return(members, nil)
		mockSession.On("ChannelMessageSend", mock.Anything, mock.Anything).Return(&discordgo.Message{}, nil)

		result, err := GetUsersWithRole(mockSession, guildID, roleID)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(result))
	})
}

func TestFormatUserMentions(t *testing.T) {
	t.Run("formats user Mentions correctly", func(t *testing.T) {
		members := []*discordgo.Member{
			{User: &discordgo.User{ID: "123"}},
			{User: &discordgo.User{ID: "456"}},
		}

		mentions := FormatUserMentions(members)
		assert.Equal(t, []string{"<@123>", "<@456>"}, mentions)
	})
	t.Run("handles empty member list", func(t *testing.T) {
		mentions := FormatUserMentions([]*discordgo.Member{})
		assert.Equal(t, []string{}, mentions)
	})
	t.Run("handles nil members list", func(t *testing.T) {
		mentions := FormatUserMentions(nil)
		assert.Equal(t, []string{}, mentions)
	})
}

func TestFormatRoleMention(t *testing.T) {
	t.Run("format roleID as mention", func(t *testing.T) {
		roleID := "123456789"
		mention := FormatRoleMention(roleID)
		assert.Equal(t, "<@&123456789>", mention)
	})
	t.Run("handles empty roleID", func(t *testing.T) {
		mention := FormatRoleMention("")
		assert.Equal(t, "<@&>", mention)
	})
}

func TestJoinMentions(t *testing.T) {
	t.Run("joins mentions with space separator", func(t *testing.T) {
		mentions := []string{"<@123>", "<@456>"}
		result := JoinMentions(mentions, " ")
		assert.Equal(t, "<@123> <@456>", result)
	})

	t.Run("joins mentions with comma separator", func(t *testing.T) {
		mentions := []string{"<@123>", "<@456>"}
		result := JoinMentions(mentions, ", ")
		assert.Equal(t, "<@123>, <@456>", result)
	})

	t.Run("handles single mention", func(t *testing.T) {
		mentions := []string{"<@123>"}
		result := JoinMentions(mentions, ", ")
		assert.Equal(t, "<@123>", result)
	})

	t.Run("handles empty mentions", func(t *testing.T) {
		result := JoinMentions([]string{}, ", ")
		assert.Equal(t, "", result)
	})

	t.Run("handles nil mentions", func(t *testing.T) {
		result := JoinMentions(nil, ", ")
		assert.Equal(t, "", result)
	})
}

func TestFormatMentionResponse(t *testing.T) {
	t.Run("formats response with message and mentions", func(t *testing.T) {
		mentions := []string{"<@123>", "<@456>"}
		message := "Hello"
		response := FormatMentionResponse(mentions, message)
		assert.Equal(t, "Hello <@123> <@456>", response)
	})

	t.Run("formats response with only mentions", func(t *testing.T) {
		mentions := []string{"<@123>", "<@456>"}
		response := FormatMentionResponse(mentions, "")
		assert.Equal(t, "<@123> <@456>", response)
	})

	t.Run("handles empty mentions", func(t *testing.T) {
		response := FormatMentionResponse([]string{}, "Hello")
		assert.Equal(t, "Sorry no user found under this role.", response)
	})

	t.Run("handles nil mentions", func(t *testing.T) {
		response := FormatMentionResponse(nil, "Hello")
		assert.Equal(t, "Sorry no user found under this role.", response)
	})
}
func TestFormatDevTitleResponse(t *testing.T) {
	roleID := "123456789"
	roleMention := "<@&" + roleID + ">"

	t.Run("formats response with no users", func(t *testing.T) {
		response := FormatDevTitleResponse([]string{}, roleID)
		expected := fmt.Sprintf("Found 0 users with the %s role", roleMention)
		assert.Equal(t, expected, response)
	})

	t.Run("formats response with single user", func(t *testing.T) {
		mentions := []string{"<@123>"}
		response := FormatDevTitleResponse(mentions, roleID)
		expected := fmt.Sprintf("Found 1 user with the %s role: %s", roleMention, mentions[0])
		assert.Equal(t, expected, response)
	})

	t.Run("formats response with multiple users", func(t *testing.T) {
		mentions := []string{"<@123>", "<@456>"}
		response := FormatDevTitleResponse(mentions, roleID)
		expected := fmt.Sprintf("Found %d users with the %s role: %s", len(mentions), roleMention, JoinMentions(mentions, ", "))
		assert.Equal(t, expected, response)
	})

	t.Run("handles nil mentions", func(t *testing.T) {
		response := FormatDevTitleResponse(nil, roleID)
		expected := fmt.Sprintf("Found 0 users with the %s role", roleMention)
		assert.Equal(t, expected, response)
	})

	t.Run("handles empty role ID", func(t *testing.T) {
		mentions := []string{"<@123>"}
		emptyRoleMention := "<@&>"
		response := FormatDevTitleResponse([]string{"<@123>"}, "")
		expected := fmt.Sprintf("Found 1 user with the %s role: %s", emptyRoleMention, mentions[0])
		assert.Equal(t, expected, response)
	})
}