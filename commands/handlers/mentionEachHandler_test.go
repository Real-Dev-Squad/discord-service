package handlers

import (
	"errors"
	"fmt"
	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
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

	panic(fmt.Sprintf("mock return value for GuildMember is not of type []*discordgo.Member: %T", args.Get(0)))
}

func (m *MockDiscordSession) ChannelMessageSend(ChannelID, content string) (*discordgo.Message, error) {
	args := m.Called(ChannelID, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	if msg, ok := args.Get(0).(*discordgo.Message); ok {
		return msg, args.Error(1)
	}

	panic(fmt.Sprintf("mock return value for ChannelMessageSend is not of type *discordgo.Message: %T", args.Get(0)))
}

func (m *MockDiscordSession) Close() error {
	args := m.Called()
	return args.Error(0)
}

var _ utils.DiscordSessionInterface = (*MockDiscordSession)(nil)

func setupTestCommandHandler(metaData map[string]string) *CommandHandler {
	return &CommandHandler{
		discordMessage: &dtos.DataPacket{
			MetaData: metaData,
		},
	}
}

func TestMentionEachHandler(t *testing.T) {
	originalCreateSession := CreateSession
	originalExtractParams := extractCommandParamsFunc
	originalFetchMembers := fetchMembersWithRoleFunc
	originalSendNoMembers := sendNoMembersMessageFunc
	originalHandleDevTitle := handleDevTitleModeFunc
	originalHandleDev := handleDevModeFunc
	originalHandleStandard := handleStandardModeFunc

	defer func() {
		CreateSession = originalCreateSession
		extractCommandParamsFunc = originalExtractParams
		fetchMembersWithRoleFunc = originalFetchMembers
		sendNoMembersMessageFunc = originalSendNoMembers
		handleDevTitleModeFunc = originalHandleDevTitle
		handleDevModeFunc = originalHandleDev
		handleStandardModeFunc = originalHandleStandard
	}()

	validMetaData := map[string]string{
		"role_id":    "testRole",
		"channel_id": "testChannel",
		"guild_id":   "testGuild",
		"message":    "Hello",
		"dev":        "false",
		"dev_title":  "false",
	}

	member1 := &discordgo.Member{User: &discordgo.User{ID: "user1"}, Roles: []string{"testRole"}}
	membersList := []*discordgo.Member{member1}
	var emptyMembersList []*discordgo.Member
	mockErr := errors.New("mock error")

	t.Run("Success Standard Mode", func(t *testing.T) {
		extractCommandParamsFunc = func(map[string]string) (CommandParams, error) {
			return CommandParams{
				RoleID: "testRole", ChannelID: "testChannel", GuildID: "testGuild", Message: "Hello", Dev: false, DevTitle: false,
			}, nil
		}
		CreateSession = func() (*discordgo.Session, error) {
			return &discordgo.Session{}, nil
		}
		fetchMembersWithRoleFunc = func(session utils.DiscordSessionInterface, guildID, roleID, channelID string) ([]*discordgo.Member, error) {
			assert.NotNil(t, session)
			assert.Equal(t, "testGuild", guildID)
			return membersList, nil
		}
		handleStandardModeFunc = func(session utils.DiscordSessionInterface, mentions []string, params CommandParams) error {
			assert.NotNil(t, session)
			assert.Contains(t, mentions, "<@user1>")
			assert.False(t, params.Dev)
			assert.False(t, params.DevTitle)
			return nil
		}
		handleDevModeFunc = func(utils.DiscordSessionInterface, []string, CommandParams) error {
			t.Fatal("handleDevMode called unexpectedly")
			return nil
		}
		handleDevTitleModeFunc = func(utils.DiscordSessionInterface, []string, CommandParams) error {
			t.Fatal("handleDevTitleMode called unexpectedly")
			return nil
		}
		sendNoMembersMessageFunc = func(utils.DiscordSessionInterface, string) error {
			t.Fatal("sendNoMembersMessage called unexpectedly")
			return nil
		}
		commandHandler := setupTestCommandHandler(validMetaData)
		err := commandHandler.mentionEachHandler()
		assert.NoError(t, err)
	})

	t.Run("Success Dev Mode", func(t *testing.T) {
		extractCommandParamsFunc = func(map[string]string) (CommandParams, error) {
			return CommandParams{
				RoleID: "testRole", ChannelID: "testChannel", GuildID: "testGuild", Message: "Dev Msg", Dev: true, DevTitle: false, // Set Dev true
			}, nil
		}
		CreateSession = func() (*discordgo.Session, error) { return &discordgo.Session{}, nil }
		fetchMembersWithRoleFunc = func(utils.DiscordSessionInterface, string, string, string) ([]*discordgo.Member, error) {
			return membersList, nil
		}
		handleDevModeFunc = func(session utils.DiscordSessionInterface, mentions []string, params CommandParams) error {
			assert.NotNil(t, session)
			assert.True(t, params.Dev)
			assert.Equal(t, "Dev Msg", params.Message)
			assert.Contains(t, mentions, "<@user1>")
			return nil
		}

		handleStandardModeFunc = func(utils.DiscordSessionInterface, []string, CommandParams) error {
			t.Fatal("handleStandardMode called unexpectedly")
			return nil
		}
		handleDevTitleModeFunc = func(utils.DiscordSessionInterface, []string, CommandParams) error {
			t.Fatal("handleDevTitleMode called unexpectedly")
			return nil
		}
		sendNoMembersMessageFunc = func(utils.DiscordSessionInterface, string) error {
			t.Fatal("sendNoMembersMessage called unexpectedly")
			return nil
		}

		commandHandler := setupTestCommandHandler(nil)
		err := commandHandler.mentionEachHandler()
		assert.NoError(t, err)
	})

	t.Run("Success DevTitle Mode", func(t *testing.T) {
		extractCommandParamsFunc = func(map[string]string) (CommandParams, error) {
			return CommandParams{
				RoleID: "testRole", ChannelID: "testChannel", GuildID: "testGuild", Message: "", Dev: false, DevTitle: true, // Set DevTitle true
			}, nil
		}
		CreateSession = func() (*discordgo.Session, error) { return &discordgo.Session{}, nil }
		fetchMembersWithRoleFunc = func(utils.DiscordSessionInterface, string, string, string) ([]*discordgo.Member, error) {
			return membersList, nil
		}
		handleDevTitleModeFunc = func(session utils.DiscordSessionInterface, mentions []string, params CommandParams) error {
			assert.NotNil(t, session)
			assert.True(t, params.DevTitle)
			assert.Contains(t, mentions, "<@user1>")
			return nil
		}
		handleStandardModeFunc = func(utils.DiscordSessionInterface, []string, CommandParams) error {
			t.Fatal("handleStandardMode called unexpectedly")
			return nil
		}
		handleDevModeFunc = func(utils.DiscordSessionInterface, []string, CommandParams) error {
			t.Fatal("handleDevMode called unexpectedly")
			return nil
		}
		sendNoMembersMessageFunc = func(utils.DiscordSessionInterface, string) error {
			t.Fatal("sendNoMembersMessage called unexpectedly")
			return nil
		}

		commandHandler := setupTestCommandHandler(nil)
		err := commandHandler.mentionEachHandler()
		assert.NoError(t, err)
	})

	t.Run("Success No Members Found", func(t *testing.T) {

		extractCommandParamsFunc = func(map[string]string) (CommandParams, error) {
			return CommandParams{ChannelID: "testChannel" /* other fields don't matter */}, nil
		}

		CreateSession = func() (*discordgo.Session, error) { return &discordgo.Session{}, nil }
		fetchMembersWithRoleFunc = func(utils.DiscordSessionInterface, string, string, string) ([]*discordgo.Member, error) {
			return emptyMembersList, nil
		}

		sendNoMembersMessageFunc = func(session utils.DiscordSessionInterface, channelID string) error {
			assert.NotNil(t, session)
			assert.Equal(t, "testChannel", channelID)
			return nil
		}

		handleStandardModeFunc = func(utils.DiscordSessionInterface, []string, CommandParams) error {
			t.Fatal("handleStandardMode called unexpectedly")
			return nil
		}

		handleDevModeFunc = func(utils.DiscordSessionInterface, []string, CommandParams) error {
			t.Fatal("handleDevMode called unexpectedly")
			return nil
		}

		handleDevTitleModeFunc = func(utils.DiscordSessionInterface, []string, CommandParams) error {
			t.Fatal("handleDevTitleMode called unexpectedly")
			return nil
		}

		commandHandler := setupTestCommandHandler(nil)
		err := commandHandler.mentionEachHandler()
		assert.NoError(t, err)
	})

	t.Run("Error ExtractCommandParams Fails", func(t *testing.T) {
		extractCommandParamsFunc = func(map[string]string) (CommandParams, error) {
			return CommandParams{}, mockErr
		}
		commandHandler := setupTestCommandHandler(nil)
		err := commandHandler.mentionEachHandler()

		assert.Error(t, err)
		assert.ErrorIs(t, err, mockErr)
		assert.ErrorContains(t, err, "failed to extract command params")
	})

	t.Run("Error CreateSession Fails", func(t *testing.T) {
		extractCommandParamsFunc = func(map[string]string) (CommandParams, error) { return CommandParams{}, nil }
		CreateSession = func() (*discordgo.Session, error) {
			return nil, mockErr
		}
		commandHandler := setupTestCommandHandler(nil)

		err := commandHandler.mentionEachHandler()
		assert.Error(t, err)
		assert.ErrorIs(t, err, mockErr)
		assert.Contains(t, err.Error(), "failed to create Discord session:")
		assert.Contains(t, err.Error(), mockErr.Error())
	})

	t.Run("Error FetchMembersWithRole Fails", func(t *testing.T) {
		extractCommandParamsFunc = func(map[string]string) (CommandParams, error) { return CommandParams{}, nil }
		CreateSession = func() (*discordgo.Session, error) { return &discordgo.Session{}, nil }
		fetchMembersWithRoleFunc = func(utils.DiscordSessionInterface, string, string, string) ([]*discordgo.Member, error) {
			return nil, mockErr
		}
		commandHandler := setupTestCommandHandler(nil)
		err := commandHandler.mentionEachHandler()
		assert.Error(t, err)
		assert.Equal(t, mockErr, err)
	})

	t.Run("Error SendNoMembersMessage Fails", func(t *testing.T) {
		extractCommandParamsFunc = func(map[string]string) (CommandParams, error) { return CommandParams{ChannelID: "testChannel"}, nil }
		CreateSession = func() (*discordgo.Session, error) { return &discordgo.Session{}, nil }
		fetchMembersWithRoleFunc = func(utils.DiscordSessionInterface, string, string, string) ([]*discordgo.Member, error) {
			return emptyMembersList, nil
		}
		sendNoMembersMessageFunc = func(utils.DiscordSessionInterface, string) error {
			return mockErr
		}
		commandHandler := setupTestCommandHandler(nil)
		err := commandHandler.mentionEachHandler()
		assert.Error(t, err)
		assert.Equal(t, mockErr, err)
	})

	t.Run("Error HandleStandardMode Fails", func(t *testing.T) {
		extractCommandParamsFunc = func(map[string]string) (CommandParams, error) { return CommandParams{Dev: false, DevTitle: false}, nil }
		CreateSession = func() (*discordgo.Session, error) { return &discordgo.Session{}, nil }
		fetchMembersWithRoleFunc = func(utils.DiscordSessionInterface, string, string, string) ([]*discordgo.Member, error) {
			return membersList, nil
		}
		handleStandardModeFunc = func(utils.DiscordSessionInterface, []string, CommandParams) error {
			return mockErr
		}
		commandHandler := setupTestCommandHandler(nil)
		err := commandHandler.mentionEachHandler()
		assert.Error(t, err)
		assert.Equal(t, mockErr, err)
	})

	t.Run("Error HandleDevMode Fails", func(t *testing.T) {
		extractCommandParamsFunc = func(map[string]string) (CommandParams, error) { return CommandParams{Dev: true, DevTitle: false}, nil }
		CreateSession = func() (*discordgo.Session, error) { return &discordgo.Session{}, nil }
		fetchMembersWithRoleFunc = func(utils.DiscordSessionInterface, string, string, string) ([]*discordgo.Member, error) {
			return membersList, nil
		}
		handleDevModeFunc = func(utils.DiscordSessionInterface, []string, CommandParams) error {
			return mockErr
		}
		commandHandler := setupTestCommandHandler(nil)

		err := commandHandler.mentionEachHandler()
		assert.Error(t, err)
		assert.Equal(t, mockErr, err)
	})

	t.Run("Error HandleDevTitleMode Fails", func(t *testing.T) {
		extractCommandParamsFunc = func(map[string]string) (CommandParams, error) { return CommandParams{Dev: false, DevTitle: true}, nil }
		CreateSession = func() (*discordgo.Session, error) { return &discordgo.Session{}, nil }
		fetchMembersWithRoleFunc = func(utils.DiscordSessionInterface, string, string, string) ([]*discordgo.Member, error) {
			return membersList, nil
		}
		handleDevTitleModeFunc = func(utils.DiscordSessionInterface, []string, CommandParams) error {
			return mockErr
		}
		commandHandler := setupTestCommandHandler(nil)
		err := commandHandler.mentionEachHandler()
		assert.Error(t, err)
		assert.Equal(t, mockErr, err)
	})

	t.Run("Success Path with session Close Error", func(t *testing.T) {
		extractCommandParamsFunc = func(map[string]string) (CommandParams, error) {
			return CommandParams{RoleID: "testRole", ChannelID: "testChannel", GuildID: "testGuild", Message: "Hello", Dev: false, DevTitle: false}, nil
		}

		mockSessionInstance := new(MockDiscordSession)
		closeErr := errors.New("failed t0 close session")
		mockSessionInstance.On("Close").Return(closeErr).Once()

		CreateSession = func() (*discordgo.Session, error) {
			return &discordgo.Session{}, nil
		}

		fetchMembersWithRoleFunc = func(session utils.DiscordSessionInterface, guildID, roleID, channelID string) ([]*discordgo.Member, error) {
			return membersList, nil
		}

		handleStandardModeFunc = func(session utils.DiscordSessionInterface, mentions []string, params CommandParams) error {
			return nil
		}

		handleDevModeFunc = func(session utils.DiscordSessionInterface, mentions []string, params CommandParams) error {
			return nil
		}
		handleDevTitleModeFunc = func(utils.DiscordSessionInterface, []string, CommandParams) error {
			t.Fatal("handleDevTitleMode called unexpectedly")
			return nil
		}
		sendNoMembersMessageFunc = func(utils.DiscordSessionInterface, string) error {
			t.Fatal("sendNoMembersMessage called unexpectedly")
			return nil
		}

		commandHandler := setupTestCommandHandler(nil)
		err := commandHandler.mentionEachHandler()
		assert.NoError(t, err)
	})
}

// commands/handlers/mentionEachHandler_test.go

// Inside TestExtractCommandParamsLogic
func TestExtractCommandParamsLogic(t *testing.T) {

	// --- Tests for Valid Cases ---
	t.Run("Valid parameters all present", func(t *testing.T) {
		metaData := map[string]string{
			"role_id": "role1", "channel_id": "chan1", "guild_id": "guild1",
			"message": "Hello", "dev": "true", "dev_title": "false",
		}
		params, err := extractCommandParamsFunc(metaData)
		assert.NoError(t, err)
		assert.Equal(t, "role1", params.RoleID)
		assert.Equal(t, "chan1", params.ChannelID)
		assert.Equal(t, "guild1", params.GuildID)
		assert.Equal(t, "Hello", params.Message)
		assert.True(t, params.Dev)
		assert.False(t, params.DevTitle)
	})

	t.Run("Valid parameters optional missing", func(t *testing.T) {
		metaData := map[string]string{
			"role_id": "role1", "channel_id": "chan1", "guild_id": "guild1",
			// message, dev, dev_title missing
		}
		params, err := extractCommandParamsFunc(metaData)
		assert.NoError(t, err)
		assert.Equal(t, "role1", params.RoleID) // Check required are still set
		assert.Equal(t, "", params.Message)     // Defaults to ""
		assert.False(t, params.Dev)             // Defaults to false
		assert.False(t, params.DevTitle)        // Defaults to false
	})

	// --- Tests for Message Length ---
	t.Run("Valid parameters with long message (truncates)", func(t *testing.T) {
		longMessage := strings.Repeat("a", MaxUserMessageLength+50)
		truncatedSuffix := "..."
		expectedMessage := strings.Repeat("a", MaxUserMessageLength) + truncatedSuffix
		metaData := map[string]string{
			"role_id": "role1", "channel_id": "chan1", "guild_id": "guild1", "message": longMessage,
		}
		params, err := extractCommandParamsFunc(metaData)
		assert.NoError(t, err)
		assert.Equal(t, expectedMessage, params.Message, "Message should be truncated")
		assert.Len(t, params.Message, MaxUserMessageLength+len(truncatedSuffix))
	})

	t.Run("Valid parameters with message within limit", func(t *testing.T) {
		shortMessage := strings.Repeat("a", MaxUserMessageLength-10)
		metaData := map[string]string{
			"role_id": "role1", "channel_id": "chan1", "guild_id": "guild1", "message": shortMessage,
		}
		params, err := extractCommandParamsFunc(metaData)
		assert.NoError(t, err)
		assert.Equal(t, shortMessage, params.Message, "Message should not be truncated")
	})

	// --- Tests for Invalid Boolean Parsing ---
	t.Run("Invalid dev parameter", func(t *testing.T) {
		metaData := map[string]string{
			"role_id": "role1", "channel_id": "chan1", "guild_id": "guild1", "dev": "not-a-bool",
		}
		params, err := extractCommandParamsFunc(metaData)
		assert.NoError(t, err) // Parsing error doesn't fail the extraction overall
		assert.False(t, params.Dev, "Dev should default to false on parse error")
	})

	t.Run("Invalid dev_title parameter", func(t *testing.T) {
		metaData := map[string]string{
			"role_id": "role1", "channel_id": "chan1", "guild_id": "guild1", "dev_title": "123",
		}
		params, err := extractCommandParamsFunc(metaData)
		assert.NoError(t, err) // Parsing error doesn't fail the extraction overall
		assert.False(t, params.DevTitle, "DevTitle should default to false on parse error")
	})

	// --- Tests for Missing/Empty Required Parameters ---
	expectedErrorMsg := "failed to extract command params: missing or empty role_id, channel_id, or guild_id" // Define once

	t.Run("Missing required parameter role_id", func(t *testing.T) {
		metaData := map[string]string{"channel_id": "chan1", "guild_id": "guild1"}
		_, err := extractCommandParamsFunc(metaData)
		assert.Error(t, err)
		// assert.Contains(t, err.Error(), "missing or empty role_id") // OLD
		assert.EqualError(t, err, expectedErrorMsg) // NEW - Check exact error string
	})

	t.Run("Empty required parameter role_id", func(t *testing.T) {
		metaData := map[string]string{"role_id": "", "channel_id": "chan1", "guild_id": "guild1"}
		_, err := extractCommandParamsFunc(metaData)
		assert.Error(t, err)
		assert.EqualError(t, err, expectedErrorMsg) // NEW
	})

	t.Run("Missing required parameter channel_id", func(t *testing.T) {
		metaData := map[string]string{"role_id": "role1", "guild_id": "guild1"}
		_, err := extractCommandParamsFunc(metaData)
		assert.Error(t, err)
		// assert.Contains(t, err.Error(), "missing or empty channel_id") // OLD
		assert.EqualError(t, err, expectedErrorMsg) // NEW
	})

	t.Run("Empty required parameter channel_id", func(t *testing.T) {
		metaData := map[string]string{"role_id": "role1", "channel_id": "", "guild_id": "guild1"}
		_, err := extractCommandParamsFunc(metaData)
		assert.Error(t, err)
		// assert.Contains(t, err.Error(), "missing or empty channel_id") // OLD
		assert.EqualError(t, err, expectedErrorMsg) // NEW
	})

	t.Run("Missing required parameter guild_id", func(t *testing.T) {
		metaData := map[string]string{"role_id": "role1", "channel_id": "chan1"}
		_, err := extractCommandParamsFunc(metaData)
		assert.Error(t, err)
		// assert.Contains(t, err.Error(), "missing or empty guild_id") // OLD
		assert.EqualError(t, err, expectedErrorMsg) // NEW
	})

	t.Run("Empty required parameter guild_id", func(t *testing.T) {
		metaData := map[string]string{"role_id": "role1", "channel_id": "chan1", "guild_id": ""}
		_, err := extractCommandParamsFunc(metaData)
		assert.Error(t, err)
		// assert.Contains(t, err.Error(), "missing or empty guild_id") // OLD
		assert.EqualError(t, err, expectedErrorMsg) // NEW
	})
}

// commands/handlers/mentionEachHandler_test.go

// --- Test logic of fetchMembersWithRoleFunc ---
func TestFetchMembersWithRoleLogic(t *testing.T) {
	// --- Test Data ---
	guildID := "g1"
	roleID := "r1"
	channelID := "c1"
	member1 := &discordgo.Member{User: &discordgo.User{ID: "u1"}, Roles: []string{roleID}}
	member2 := &discordgo.Member{User: &discordgo.User{ID: "u2"}, Roles: []string{"other"}}
	membersListInputSuccess := []*discordgo.Member{member1, member2} // Input for success case mock
	expectedOutputSuccess := []*discordgo.Member{member1}            // Expected result after filtering
	mockErr := errors.New("API error from GuildMembers")             // Error returned BY GuildMembers mock
	var emptyMemberList []*discordgo.Member

	t.Run("Success - Members found", func(t *testing.T) {
		mockSession := new(MockDiscordSession)

		// --- Mock Expectations for session calls ---
		// 1. Expect first call to GuildMembers (via utils.GetUsersWithRole)
		mockSession.On("GuildMembers", guildID, "", 1000).Return(membersListInputSuccess, nil).Once()
		// 2. Expect second call after processing page 1 (last ID "u2") -> return empty
		mockSession.On("GuildMembers", guildID, member2.User.ID, 1000).Return(emptyMemberList, nil).Once()
		// --- End Mock Expectations ---

		// Act: Call the REAL fetchMembersWithRoleFunc variable
		members, err := fetchMembersWithRoleFunc(mockSession, guildID, roleID, channelID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedOutputSuccess, members) // Check filtering result
		mockSession.AssertExpectations(t)               // Verify both GuildMembers calls happened
		// ChannelMessageSend should NOT be called on success path
		mockSession.AssertNotCalled(t, "ChannelMessageSend", mock.Anything, mock.Anything)
	})

	t.Run("Success - No members with role found", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		membersInputNoMatch := []*discordgo.Member{member2} // Input has no matching role

		// --- Mock Expectations ---
		// 1. Expect first call
		mockSession.On("GuildMembers", guildID, "", 1000).Return(membersInputNoMatch, nil).Once()
		// 2. Expect second call after processing page 1 (last ID "u2") -> return empty
		mockSession.On("GuildMembers", guildID, member2.User.ID, 1000).Return(emptyMemberList, nil).Once()
		// --- End Mock Expectations ---

		// Act
		members, err := fetchMembersWithRoleFunc(mockSession, guildID, roleID, channelID)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, members)          // Result should be empty after filtering
		mockSession.AssertExpectations(t) // Verify both GuildMembers calls happened
		mockSession.AssertNotCalled(t, "ChannelMessageSend", mock.Anything, mock.Anything)
	})

	t.Run("Error - GuildMembers fails, sending error message succeeds", func(t *testing.T) {
		mockSession := new(MockDiscordSession)

		// --- Mock the underlying session call to fail ---
		mockSession.On("GuildMembers", guildID, "", 1000).Return(nil, mockErr).Once() // Simulate GuildMembers failing

		// --- Expect ChannelMessageSend to be called DIRECTLY by fetchMembersWithRoleFunc ---
		expectedErrorMsg := fmt.Sprintf("Sorry, I couldn't fetch members for role <@&%s> right now. Please try again later.", roleID)
		mockSession.On("ChannelMessageSend", channelID, expectedErrorMsg).Return(&discordgo.Message{}, nil).Once()
		// --- End Expect ---

		// Act: Call the REAL fetchMembersWithRoleFunc
		members, err := fetchMembersWithRoleFunc(mockSession, guildID, roleID, channelID)

		// Assert
		assert.Error(t, err)
		// The error returned by fetchMembersWithRoleFunc is the wrapped one from utils.GetUsersWithRole
		assert.ErrorContains(t, err, "failed to fetch guild members chunk:")
		assert.ErrorContains(t, err, mockErr.Error())
		assert.Nil(t, members)
		mockSession.AssertExpectations(t) // Verify BOTH GuildMembers and ChannelMessageSend were called
	})

	t.Run("Error - GuildMembers fails, sending error message fails", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		sendErr := errors.New("send failed")

		// --- Mock the underlying session call to fail ---
		mockSession.On("GuildMembers", guildID, "", 1000).Return(nil, mockErr).Once() // Simulate GuildMembers failing

		// --- Expect ChannelMessageSend call, but mock its return to fail ---
		expectedErrorMsg := fmt.Sprintf("Sorry, I couldn't fetch members for role <@&%s> right now. Please try again later.", roleID)
		mockSession.On("ChannelMessageSend", channelID, expectedErrorMsg).Return(nil, sendErr).Once()
		// --- End Expect ---

		// Act: Call the REAL fetchMembersWithRoleFunc
		members, err := fetchMembersWithRoleFunc(mockSession, guildID, roleID, channelID)

		// Assert
		assert.Error(t, err)
		// Still expect the original fetch error (the wrapped one from GetUsersWithRole) to be returned
		assert.ErrorContains(t, err, "failed to fetch guild members chunk:")
		assert.ErrorContains(t, err, mockErr.Error())
		assert.Nil(t, members)
		mockSession.AssertExpectations(t) // Verify BOTH GuildMembers and ChannelMessageSend were called
	})
}

// ... rest of file ...

func TestSendNoMembersMessageLogic(t *testing.T) {
	channelID := "c1"
	expectedMsg := "Sorry, no members found with this role"
	mockErr := errors.New("send failed")

	t.Run("Success sending message", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		mockSession.On("ChannelMessageSend", channelID, expectedMsg).Return(&discordgo.Message{}, nil).Once()

		err := sendNoMembersMessageFunc(mockSession, channelID)

		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Error sending message", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		mockSession.On("ChannelMessageSend", channelID, expectedMsg).Return(nil, mockErr).Once()

		err := sendNoMembersMessageFunc(mockSession, channelID)

		assert.Error(t, err)
		assert.Equal(t, mockErr, err)
		mockSession.AssertExpectations(t)
	})
}

func TestHandleDevModeLogic(t *testing.T) {
	mockErr := errors.New("send failed")
	params := CommandParams{ChannelID: "c1", Message: "Test", RoleID: "r1"}
	mentions := []string{"<@u1>", "<@u2>", "<@u3>", "<@u4>", "<@u5>", "<@u6>"} // More than BatchSize

	t.Run("Success with batching", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		for i := 0; i < 6; i++ {
			expectedMsg := fmt.Sprintf("%s %s", params.Message, mentions[i])
			mockSession.On("ChannelMessageSend", params.ChannelID, expectedMsg).Return(&discordgo.Message{}, nil).Once()
		}

		err := handleDevModeFunc(mockSession, mentions, params)

		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Success single batch", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		shortMentions := []string{"<@uA>", "<@uB>"}
		mockSession.On("ChannelMessageSend", params.ChannelID, fmt.Sprintf("%s %s", params.Message, shortMentions[0])).Return(&discordgo.Message{}, nil).Once()
		mockSession.On("ChannelMessageSend", params.ChannelID, fmt.Sprintf("%s %s", params.Message, shortMentions[1])).Return(&discordgo.Message{}, nil).Once()

		err := handleDevModeFunc(mockSession, shortMentions, params)

		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Success no mentions", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		err := handleDevModeFunc(mockSession, []string{}, params)
		assert.NoError(t, err)
		mockSession.AssertNotCalled(t, "ChannelMessageSend", mock.Anything, mock.Anything)
	})

	t.Run("Success no custom message", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		localParams := CommandParams{ChannelID: "c1", Message: "", RoleID: "r1"}
		mention := "<@uOnly>"
		mockSession.On("ChannelMessageSend", localParams.ChannelID, mention).Return(&discordgo.Message{}, nil).Once()

		err := handleDevModeFunc(mockSession, []string{mention}, localParams)

		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Error during sending some mentions", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		shortMentions := []string{"<@uA>", "<@uB>", "<@uC>"}
		mockSession.On("ChannelMessageSend", params.ChannelID, fmt.Sprintf("%s %s", params.Message, shortMentions[0])).Return(&discordgo.Message{}, nil).Once()
		mockSession.On("ChannelMessageSend", params.ChannelID, fmt.Sprintf("%s %s", params.Message, shortMentions[1])).Return(nil, mockErr).Once() // Fail here
		mockSession.On("ChannelMessageSend", params.ChannelID, fmt.Sprintf("%s %s", params.Message, shortMentions[2])).Return(&discordgo.Message{}, nil).Once()

		expectedSummaryMsg := fmt.Sprintf("Finished mentioning, but failed for 1 users: %s", shortMentions[1])
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedSummaryMsg).Return(&discordgo.Message{}, nil).Once() // Expect summary message

		err := handleDevModeFunc(mockSession, shortMentions, params)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to send 1 out of 3 mentions")
		mockSession.AssertExpectations(t)
	})

}

func TestHandleUserListModeLogic(t *testing.T) {
	mockErr := errors.New("send failed")
	params := CommandParams{ChannelID: "c1", RoleID: "r1"}
	mentions := []string{"<@u1>", "<@u2>"}
	expectedResponse := utils.FormatUserListResponse(mentions, params.RoleID)

	t.Run("Success sending message", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedResponse).Return(&discordgo.Message{}, nil).Once()

		err := handleDevTitleModeFunc(mockSession, mentions, params)

		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Error sending message", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedResponse).Return(nil, mockErr).Once()

		err := handleDevTitleModeFunc(mockSession, mentions, params)

		assert.Error(t, err)
		assert.Equal(t, mockErr, err)
		mockSession.AssertExpectations(t)
	})
}

func TestHandleStandardModeLogic(t *testing.T) {
	mockErr := errors.New("send failed")
	params := CommandParams{ChannelID: "c1", RoleID: "r1", Message: "Hi"}
	mentions := []string{"<@u1>", "<@u2>"}
	expectedResponse := utils.FormatMentionResponse(mentions, params.Message)

	t.Run("Success sending message", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedResponse).Return(&discordgo.Message{}, nil).Once()

		err := handleStandardModeFunc(mockSession, mentions, params)

		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Error sending message", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedResponse).Return(nil, mockErr).Once()

		err := handleStandardModeFunc(mockSession, mentions, params)

		assert.Error(t, err)
		assert.Equal(t, mockErr, err)
		mockSession.AssertExpectations(t)
	})
}
