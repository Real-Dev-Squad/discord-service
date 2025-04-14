package utils

import (
	"errors"
	"fmt"
	"strings"
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if members, ok := args.Get(0).([]*discordgo.Member); ok {
		return members, args.Error(1)
	}
	panic(fmt.Sprintf("mock return value for GuildMembers is not []*discordgo.Member: %T", args.Get(0)))
}

func (m *MockDiscordSession) ChannelMessageSend(channelID, content string) (*discordgo.Message, error) {
	args := m.Called(channelID, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if msg, ok := args.Get(0).(*discordgo.Message); ok {
		return msg, args.Error(1)
	}
	panic(fmt.Sprintf("mock return value for ChannelMessageSend is not *discordgo.Message: %T", args.Get(0)))
}

// TestGetUsersWithRole tests the GetUsersWithRole function which is responsible for
// fetching members with a specific role, handling pagination implicitly via the session interface.
// It uses MockDiscordSession to simulate responses from the Discord API (GuildMembers call).
func TestGetUsersWithRole(t *testing.T) {
	guildID := "testGuild"
	roleID := "testRole"

	member1 := &discordgo.Member{User: &discordgo.User{ID: "123"}, Roles: []string{roleID}}
	member2 := &discordgo.Member{User: &discordgo.User{ID: "456"}, Roles: []string{"otherRole"}}
	member3 := &discordgo.Member{User: &discordgo.User{ID: "789"}, Roles: []string{roleID, "anotherRole"}}

	var memberNilUser *discordgo.Member = &discordgo.Member{User: nil, Roles: []string{roleID}}
	var memberNilRoles *discordgo.Member = &discordgo.Member{User: &discordgo.User{ID: "abc"}, Roles: nil}
	var memberNil *discordgo.Member = nil
	var emptyMemberList []*discordgo.Member

	t.Run("returns single user with matching role", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		membersInput := []*discordgo.Member{member1, member2}
		mockSession.On("GuildMembers", guildID, "", 1000).Return(membersInput, nil).Once()

		result, err := GetUsersWithRole(mockSession, guildID, roleID)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "123", result[0].User.ID)
		mockSession.AssertExpectations(t)
	})

	t.Run("returns multiple users with matching role", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		membersInput := []*discordgo.Member{member1, member2, member3}
		mockSession.On("GuildMembers", guildID, "", 1000).Return(membersInput, nil).Once()

		result, err := GetUsersWithRole(mockSession, guildID, roleID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		expectedIDs := []string{"123", "789"}
		var actualIDs []string
		for _, m := range result {
			if m != nil && m.User != nil {
				actualIDs = append(actualIDs, m.User.ID)
			}
		}
		assert.ElementsMatch(t, expectedIDs, actualIDs, "Should find members 123 and 789 only")
		mockSession.AssertExpectations(t)
	})

	t.Run("handles error from GuildMembers", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		mockErr := errors.New("API error")
		mockSession.On("GuildMembers", guildID, "", 1000).Return(nil, mockErr).Once()

		_, err := GetUsersWithRole(mockSession, guildID, roleID)

		assert.Error(t, err)
		assert.ErrorContains(t, err, mockErr.Error())
		mockSession.AssertExpectations(t)
	})

	t.Run("returns empty slice when no users have the role", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		membersInput := []*discordgo.Member{member2}
		mockSession.On("GuildMembers", guildID, "", 1000).Return(membersInput, nil).Once()

		result, err := GetUsersWithRole(mockSession, guildID, roleID)

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockSession.AssertExpectations(t)
	})

	t.Run("handles empty member list from GuildMembers", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		mockSession.On("GuildMembers", guildID, "", 1000).Return(emptyMemberList, nil).Once() // Use var
		result, err := GetUsersWithRole(mockSession, guildID, roleID)
		assert.NoError(t, err)
		assert.Empty(t, result)
		mockSession.AssertExpectations(t)
	})

	t.Run("skips members with nil User or nil Roles safely", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		membersInput := []*discordgo.Member{member1, memberNilUser, memberNilRoles, memberNil, member3}
		mockSession.On("GuildMembers", guildID, "", 1000).Return(membersInput, nil).Once()
		result, err := GetUsersWithRole(mockSession, guildID, roleID)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		expectedIDs := []string{"123", "789"}
		var actualIDs []string
		for _, m := range result {
			if m != nil && m.User != nil {
				actualIDs = append(actualIDs, m.User.ID)
			}
		}
		assert.ElementsMatch(t, expectedIDs, actualIDs)
		mockSession.AssertExpectations(t)
	})
}

// TestFormatUserMentions tests the utility function for converting member objects
// into Discord mention strings.
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
		assert.Empty(t, mentions)
	})
	t.Run("handles nil members list", func(t *testing.T) {
		mentions := FormatUserMentions(nil)
		assert.Empty(t, mentions)
	})
}

// TestFormatRoleMention tests the utility function for formatting role IDs
// into Discord role mention strings.
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

// TestFormatMentionResponse tests the utility function for creating the final message
// content for the standard mention-each mode.
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
		roleID := "123456789"
		roleMention := "<@&" + roleID + ">"
		mentions := []string{"<@123>", "<@456>"}
		response := FormatDevTitleResponse(mentions, roleID)
		expected := fmt.Sprintf("Found %d users with the %s role: %s", len(mentions), roleMention, strings.Join(mentions, ", "))
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
