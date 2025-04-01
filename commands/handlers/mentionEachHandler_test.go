package handlers

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDiscordSession duplicates the mock from members_utils_test.go for locality
type MockDiscordSession struct {
	mock.Mock
}

func (m *MockDiscordSession) GuildMembers(guildID, after string, limit int) ([]*discordgo.Member, error) {
	args := m.Called(guildID, after, limit)
	// Handle potential nil return for the member slice
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*discordgo.Member), args.Error(1)
}

func (m *MockDiscordSession) ChannelMessageSend(channelID, content string) (*discordgo.Message, error) {
	args := m.Called(channelID, content)
	// Handle potential nil return for the message
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*discordgo.Message), args.Error(1)
}

// Ensure MockDiscordSession implements the interface
var _ utils.DiscordSessionInterface = (*MockDiscordSession)(nil)

func TestExtractCommandParams(t *testing.T) {
	t.Run("Valid parameters", func(t *testing.T) {
		metaData := map[string]string{
			"role_id":    "testRole",
			"channel_id": "testChannel",
			"guild_id":   "testGuild",
			"message":    "Hello",
			"dev":        "true",
			"dev_title":  "false",
		}
		params, err := extractCommandParams(metaData)
		assert.NoError(t, err)
		assert.Equal(t, "testRole", params.RoleID)
		assert.Equal(t, "testChannel", params.ChannelID)
		assert.Equal(t, "testGuild", params.GuildID)
		assert.Equal(t, "Hello", params.Message)
		assert.True(t, params.Dev)
		assert.False(t, params.DevTitle)
	})

	t.Run("Missing required parameter role_id", func(t *testing.T) {
		metaData := map[string]string{
			"channel_id": "testChannel",
			"guild_id":   "testGuild",
		}
		_, err := extractCommandParams(metaData)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to extract command params")
	})

	t.Run("Missing required parameter channel_id", func(t *testing.T) {
		metaData := map[string]string{
			"role_id":  "testRole",
			"guild_id": "testGuild",
		}
		_, err := extractCommandParams(metaData)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to extract command params")
	})

	t.Run("Missing required parameter guild_id", func(t *testing.T) {
		metaData := map[string]string{
			"role_id":    "testRole",
			"channel_id": "testChannel",
		}
		_, err := extractCommandParams(metaData)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to extract command params")
	})

	t.Run("Optional parameters missing", func(t *testing.T) {
		metaData := map[string]string{
			"role_id":    "testRole",
			"channel_id": "testChannel",
			"guild_id":   "testGuild",
		}
		params, err := extractCommandParams(metaData)
		assert.NoError(t, err)
		assert.Equal(t, "", params.Message)
		assert.False(t, params.Dev)
		assert.False(t, params.DevTitle)
	})

	t.Run("Invalid dev parameter", func(t *testing.T) {
		metaData := map[string]string{
			"role_id":    "testRole",
			"channel_id": "testChannel",
			"guild_id":   "testGuild",
			"dev":        "not-a-bool",
		}
		params, err := extractCommandParams(metaData)
		assert.NoError(t, err)
		assert.False(t, params.Dev, "Dev should default to false on parse error")
	})

	t.Run("Invalid dev_title parameter", func(t *testing.T) {
		metaData := map[string]string{
			"role_id":    "testRole",
			"channel_id": "testChannel",
			"guild_id":   "testGuild",
			"dev_title":  "not-a-bool",
		}
		params, err := extractCommandParams(metaData)
		assert.NoError(t, err)
		assert.False(t, params.DevTitle, "DevTitle should default to false on parse error")
	})
}

func TestFetchMembersWithRole(t *testing.T) {
	guildID := "testGuild"
	roleID := "testRole"
	channelID := "testChannel" // Used for sending error messages

	member1 := &discordgo.Member{User: &discordgo.User{ID: "user1"}, Roles: []string{roleID}}
	member2 := &discordgo.Member{User: &discordgo.User{ID: "user2"}, Roles: []string{"otherRole"}}
	member3 := &discordgo.Member{User: &discordgo.User{ID: "user3"}, Roles: []string{roleID, "anotherRole"}}

	t.Run("Success - Members found", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		allMembers := []*discordgo.Member{member1, member2, member3}
		mockSession.On("GuildMembers", guildID, "", 1000).Return(allMembers, nil).Once()

		members, err := fetchMembersWithRole(mockSession, guildID, roleID, channelID)
		assert.NoError(t, err)
		assert.Len(t, members, 2)
		assert.Contains(t, members, member1)
		assert.Contains(t, members, member3)
		assert.NotContains(t, members, member2)
		mockSession.AssertExpectations(t)
	})

	t.Run("Success - No members with role", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		allMembers := []*discordgo.Member{member2} // Only member with other role
		mockSession.On("GuildMembers", guildID, "", 1000).Return(allMembers, nil).Once()

		members, err := fetchMembersWithRole(mockSession, guildID, roleID, channelID)
		assert.NoError(t, err)
		assert.Empty(t, members)
		mockSession.AssertExpectations(t)
	})

	t.Run("Success - No members in guild", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		var allMembers []*discordgo.Member
		mockSession.On("GuildMembers", guildID, "", 1000).Return(allMembers, nil).Once()

		members, err := fetchMembersWithRole(mockSession, guildID, roleID, channelID)
		assert.NoError(t, err)
		assert.Empty(t, members)
		mockSession.AssertExpectations(t)
	})

	t.Run("Error fetching members", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		apiError := errors.New("discord API error")
		mockSession.On("GuildMembers", guildID, "", 1000).Return(nil, apiError).Once()

		// Expect error message to be sent to the channel
		expectedErrorMsg := fmt.Sprintf("Failed to fetch members with role: %v", fmt.Errorf("failed to fetch guild members: %w", apiError))
		mockSession.On("ChannelMessageSend", channelID, expectedErrorMsg).Return(&discordgo.Message{}, nil).Once()

		members, err := fetchMembersWithRole(mockSession, guildID, roleID, channelID)
		assert.Error(t, err)
		assert.Nil(t, members)
		assert.Contains(t, err.Error(), "discord API error")
		mockSession.AssertExpectations(t)
	})

	t.Run("Error fetching members and error sending message", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		apiError := errors.New("discord API error")
		sendError := errors.New("cannot send message")
		mockSession.On("GuildMembers", guildID, "", 1000).Return(nil, apiError).Once()

		expectedErrorMsg := fmt.Sprintf("Failed to fetch members with role: %v", fmt.Errorf("failed to fetch guild members: %w", apiError))
		mockSession.On("ChannelMessageSend", channelID, expectedErrorMsg).Return(nil, sendError).Once()

		members, err := fetchMembersWithRole(mockSession, guildID, roleID, channelID)
		assert.Error(t, err) // The original error from GuildMembers should be returned
		assert.Nil(t, members)
		assert.Contains(t, err.Error(), "discord API error")
		mockSession.AssertExpectations(t) // Verify ChannelMessageSend was still called
	})
}

func TestSendNoMembersMessage(t *testing.T) {
	channelID := "testChannel"

	t.Run("Success sending message", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		expectedMsg := "Sorry, no members found with this role"
		mockSession.On("ChannelMessageSend", channelID, expectedMsg).Return(&discordgo.Message{}, nil).Once()

		err := sendNoMembersMessage(mockSession, channelID)
		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Error sending message", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		sendError := errors.New("cannot send message")
		expectedMsg := "Sorry, no members found with this role"
		mockSession.On("ChannelMessageSend", channelID, expectedMsg).Return(nil, sendError).Once()

		err := sendNoMembersMessage(mockSession, channelID)
		assert.Error(t, err)
		assert.Equal(t, sendError, err)
		mockSession.AssertExpectations(t)
	})
}

func TestHandleDevMode(t *testing.T) {
	params := CommandParams{
		ChannelID: "testChannel",
		Message:   "Test message",
	}
	// Example: 7 mentions, BatchSize 5 -> 2 batches (5, 2)
	mentions := make([]string, 7)
	for i := 0; i < 7; i++ {
		mentions[i] = fmt.Sprintf("<@user%d>", i)
	}

	t.Run("Success with batching", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		// Batch 1 (5 messages)
		for i := 0; i < 5; i++ {
			expectedMsg := fmt.Sprintf("%s %s", params.Message, mentions[i])
			mockSession.On("ChannelMessageSend", params.ChannelID, expectedMsg).Return(&discordgo.Message{}, nil).Once()
		}
		// Batch 2 (2 messages)
		for i := 5; i < 7; i++ {
			expectedMsg := fmt.Sprintf("%s %s", params.Message, mentions[i])
			mockSession.On("ChannelMessageSend", params.ChannelID, expectedMsg).Return(&discordgo.Message{}, nil).Once()
		}

		// We can't easily test the time.Sleep duration, but we can run the function
		startTime := time.Now()
		err := handleDevMode(mockSession, mentions, params)
		duration := time.Since(startTime)

		assert.NoError(t, err)
		// Check if roughly the expected delay occurred (allowing for execution time)
		// This assertion is flaky but demonstrates intent.
		assert.GreaterOrEqual(t, duration, BatchDelay)
		assert.Less(t, duration, BatchDelay*2) // Should only sleep once

		mockSession.AssertExpectations(t)
	})

	t.Run("Success single batch", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		shortMentions := []string{"<@userA>", "<@userB>"} // Less than BatchSize

		for _, mention := range shortMentions {
			expectedMsg := fmt.Sprintf("%s %s", params.Message, mention)
			mockSession.On("ChannelMessageSend", params.ChannelID, expectedMsg).Return(&discordgo.Message{}, nil).Once()
		}

		startTime := time.Now()
		err := handleDevMode(mockSession, shortMentions, params)
		duration := time.Since(startTime)

		assert.NoError(t, err)
		assert.Less(t, duration, BatchDelay, "Should not sleep if only one batch")
		mockSession.AssertExpectations(t)
	})

	t.Run("Success no mentions", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		var emptyMentions []string

		err := handleDevMode(mockSession, emptyMentions, params)
		assert.NoError(t, err)
		mockSession.AssertExpectations(t) // No calls expected
		mockSession.AssertNotCalled(t, "ChannelMessageSend", mock.Anything, mock.Anything)
	})

	t.Run("Success no custom message", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		localParams := CommandParams{ChannelID: "testChannel", Message: ""}
		mention := "<@userOnly>"
		mockSession.On("ChannelMessageSend", localParams.ChannelID, mention).Return(&discordgo.Message{}, nil).Once()

		err := handleDevMode(mockSession, []string{mention}, localParams)
		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Error during sending", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		sendError := errors.New("cannot send")

		// Expect first message to be attempted
		expectedMsg1 := fmt.Sprintf("%s %s", params.Message, mentions[0])
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedMsg1).Return(nil, sendError).Once()

		// Don't expect subsequent messages
		// expectedMsg2 := fmt.Sprintf("%s %s", params.Message, mentions[1])
		// mockSession.AssertNotCalled(t, "ChannelMessageSend", params.ChannelID, expectedMsg2)

		err := handleDevMode(mockSession, mentions, params)
		assert.Error(t, err)
		assert.Equal(t, sendError, err)
		mockSession.AssertExpectations(t)
	})
}

func TestHandleDevTitleMode(t *testing.T) {
	params := CommandParams{
		ChannelID: "testChannel",
		RoleID:    "testRole",
	}
	mentions := []string{"<@user1>", "<@user2>"}

	t.Run("Success multiple users", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		expectedResponse := utils.FormatDevTitleResponse(mentions, params.RoleID) // Use the actual util function
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedResponse).Return(&discordgo.Message{}, nil).Once()

		err := handleDevTitleMode(mockSession, mentions, params)
		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Success single user", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		singleMention := []string{"<@user1>"}
		expectedResponse := utils.FormatDevTitleResponse(singleMention, params.RoleID)
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedResponse).Return(&discordgo.Message{}, nil).Once()

		err := handleDevTitleMode(mockSession, singleMention, params)
		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Success no users", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		var noMentions []string
		expectedResponse := utils.FormatDevTitleResponse(noMentions, params.RoleID)
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedResponse).Return(&discordgo.Message{}, nil).Once()

		err := handleDevTitleMode(mockSession, noMentions, params)
		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Error sending message", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		sendError := errors.New("cannot send")
		expectedResponse := utils.FormatDevTitleResponse(mentions, params.RoleID)
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedResponse).Return(nil, sendError).Once()

		err := handleDevTitleMode(mockSession, mentions, params)
		assert.Error(t, err)
		assert.Equal(t, sendError, err)
		mockSession.AssertExpectations(t)
	})
}

func TestHandleStandardMode(t *testing.T) {
	params := CommandParams{
		ChannelID: "testChannel",
		Message:   "Hello there",
		RoleID:    "testRole", // Used in logging
	}
	mentions := []string{"<@user1>", "<@user2>"}

	t.Run("Success with message", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		expectedResponse := utils.FormatMentionResponse(mentions, params.Message) // Use the actual util function
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedResponse).Return(&discordgo.Message{}, nil).Once()

		err := handleStandardMode(mockSession, mentions, params)
		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Success without message", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		localParams := CommandParams{ChannelID: "testChannel", Message: "", RoleID: "testRole"}
		expectedResponse := utils.FormatMentionResponse(mentions, localParams.Message)
		mockSession.On("ChannelMessageSend", localParams.ChannelID, expectedResponse).Return(&discordgo.Message{}, nil).Once()

		err := handleStandardMode(mockSession, mentions, localParams)
		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Success no users", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		var noMentions []string
		expectedResponse := utils.FormatMentionResponse(noMentions, params.Message) // Should return "Sorry no user..."
		assert.Contains(t, expectedResponse, "Sorry no user")                       // Check content
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedResponse).Return(&discordgo.Message{}, nil).Once()

		err := handleStandardMode(mockSession, noMentions, params)
		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Error sending message", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		sendError := errors.New("cannot send")
		expectedResponse := utils.FormatMentionResponse(mentions, params.Message)
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedResponse).Return(nil, sendError).Once()

		err := handleStandardMode(mockSession, mentions, params)
		assert.Error(t, err)
		assert.Equal(t, sendError, err)
		mockSession.AssertExpectations(t)
	})
}

// --- Integration-Style Test for mentionEachHandler (Focus on helper orchestration) ---
// Note: This doesn't mock CreateSession or Close, focusing on the flow after session setup.
// Direct testing of mentionEachHandler requires dependency injection for the session.
// However, testing the helpers thoroughly provides high confidence.

// Example: Simulating a Standard Mode flow via helpers
func TestMentionEachHandlerFlow_StandardMode(t *testing.T) {
	// 1. Setup Mocks & Data
	mockSession := new(MockDiscordSession)
	guildID := "guild1"
	channelID := "channel1"
	roleID := "role1"
	message := "Hey folks"
	metaData := map[string]string{
		"guild_id":   guildID,
		"channel_id": channelID,
		"role_id":    roleID,
		"message":    message,
		"dev":        "false", // Explicitly standard mode
		"dev_title":  "false",
	}

	member1 := &discordgo.Member{User: &discordgo.User{ID: "user1"}, Roles: []string{roleID}}
	member2 := &discordgo.Member{User: &discordgo.User{ID: "user2"}, Roles: []string{roleID}}
	allMembers := []*discordgo.Member{member1, member2}
	mentions := []string{"<@user1>", "<@user2>"}

	// 2. Define Expected Calls (Reverse order of execution often helps)
	//    - Expect final message send from handleStandardMode
	expectedFinalMsg := utils.FormatMentionResponse(mentions, message)
	mockSession.On("ChannelMessageSend", channelID, expectedFinalMsg).Return(&discordgo.Message{}, nil).Once()
	//    - Expect GuildMembers call from fetchMembersWithRole
	mockSession.On("GuildMembers", guildID, "", 1000).Return(allMembers, nil).Once()

	// 3. Execute parts of the flow (using the tested helpers)
	//    - Simulate param extraction
	params, errExtract := extractCommandParams(metaData)
	assert.NoError(t, errExtract)
	assert.Equal(t, guildID, params.GuildID) // Verify extraction worked as expected
	assert.Equal(t, roleID, params.RoleID)
	assert.Equal(t, channelID, params.ChannelID)
	assert.Equal(t, message, params.Message)
	assert.False(t, params.Dev)
	assert.False(t, params.DevTitle)

	//    - Simulate member fetching
	fetchedMembers, errFetch := fetchMembersWithRole(mockSession, params.GuildID, params.RoleID, params.ChannelID)
	assert.NoError(t, errFetch)
	assert.Len(t, fetchedMembers, 2)

	//    - Check if members were found (skip no members message)
	assert.NotEmpty(t, fetchedMembers) // Simulate check in main handler

	//    - Format mentions (using the util directly as the handler would)
	formattedMentions := utils.FormatUserMentions(fetchedMembers)
	assert.Equal(t, mentions, formattedMentions)

	//    - Simulate mode handling (call the specific handler)
	errHandle := handleStandardMode(mockSession, formattedMentions, params)
	assert.NoError(t, errHandle)

	// 4. Assert Mock Expectations
	mockSession.AssertExpectations(t)
}

// Example: Simulating a Dev Mode flow via helpers
func TestMentionEachHandlerFlow_DevMode(t *testing.T) {
	// 1. Setup Mocks & Data
	mockSession := new(MockDiscordSession)
	guildID := "guild1"
	channelID := "channel1"
	roleID := "role1"
	message := "Dev ping"
	metaData := map[string]string{
		"guild_id":   guildID,
		"channel_id": channelID,
		"role_id":    roleID,
		"message":    message,
		"dev":        "true", // Dev mode
		"dev_title":  "false",
	}

	// Use fewer members to simplify testing dev mode calls (e.g., 2 members, BatchSize 5 -> 1 batch)
	member1 := &discordgo.Member{User: &discordgo.User{ID: "user1"}, Roles: []string{roleID}}
	member2 := &discordgo.Member{User: &discordgo.User{ID: "user2"}, Roles: []string{roleID}}
	allMembers := []*discordgo.Member{member1, member2}
	mentions := []string{"<@user1>", "<@user2>"}

	// 2. Define Expected Calls
	//    - Expect individual message sends from handleDevMode
	expectedMsg1 := fmt.Sprintf("%s %s", message, mentions[0])
	expectedMsg2 := fmt.Sprintf("%s %s", message, mentions[1])
	mockSession.On("ChannelMessageSend", channelID, expectedMsg1).Return(&discordgo.Message{}, nil).Once()
	mockSession.On("ChannelMessageSend", channelID, expectedMsg2).Return(&discordgo.Message{}, nil).Once()
	//    - Expect GuildMembers call from fetchMembersWithRole
	mockSession.On("GuildMembers", guildID, "", 1000).Return(allMembers, nil).Once()

	// 3. Execute parts of the flow
	params, errExtract := extractCommandParams(metaData)
	assert.NoError(t, errExtract)
	assert.True(t, params.Dev) // Verify dev flag parsed

	fetchedMembers, errFetch := fetchMembersWithRole(mockSession, params.GuildID, params.RoleID, params.ChannelID)
	assert.NoError(t, errFetch)
	assert.Len(t, fetchedMembers, 2)

	assert.NotEmpty(t, fetchedMembers)

	formattedMentions := utils.FormatUserMentions(fetchedMembers)
	assert.Equal(t, mentions, formattedMentions)

	// Simulate mode handling
	errHandle := handleDevMode(mockSession, formattedMentions, params)
	assert.NoError(t, errHandle)

	// 4. Assert Mock Expectations
	mockSession.AssertExpectations(t)
}

// Example: Simulating No Members Found flow
func TestMentionEachHandlerFlow_NoMembers(t *testing.T) {
	// 1. Setup Mocks & Data
	mockSession := new(MockDiscordSession)
	guildID := "guild1"
	channelID := "channel1"
	roleID := "role1"
	metaData := map[string]string{
		"guild_id":   guildID,
		"channel_id": channelID,
		"role_id":    roleID,
		"message":    "",
		"dev":        "false",
		"dev_title":  "false",
	}

	var allMembers []*discordgo.Member // No members returned

	// 2. Define Expected Calls
	//    - Expect "no members" message send from sendNoMembersMessage
	expectedNoMembersMsg := "Sorry, no members found with this role"
	mockSession.On("ChannelMessageSend", channelID, expectedNoMembersMsg).Return(&discordgo.Message{}, nil).Once()
	//    - Expect GuildMembers call from fetchMembersWithRole
	mockSession.On("GuildMembers", guildID, "", 1000).Return(allMembers, nil).Once()

	// 3. Execute parts of the flow
	params, errExtract := extractCommandParams(metaData)
	assert.NoError(t, errExtract)

	fetchedMembers, errFetch := fetchMembersWithRole(mockSession, params.GuildID, params.RoleID, params.ChannelID)
	assert.NoError(t, errFetch)
	assert.Empty(t, fetchedMembers) // Verify fetch returned empty

	// Simulate the check in the main handler `if len(members) == 0`
	if len(fetchedMembers) == 0 {
		errSend := sendNoMembersMessage(mockSession, params.ChannelID)
		assert.NoError(t, errSend)
	} else {
		t.Errorf("Flow should have triggered sendNoMembersMessage")
	}

	// 4. Assert Mock Expectations
	mockSession.AssertExpectations(t)
}
