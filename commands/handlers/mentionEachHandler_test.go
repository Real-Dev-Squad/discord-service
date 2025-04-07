package handlers

import (
	"errors"
	"fmt"
	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestExtractCommandParamsLogic(t *testing.T) {
	t.Run("Valid parameters", func(t *testing.T) {
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

	t.Run("Valid parameters dev_title true", func(t *testing.T) {
		metaData := map[string]string{
			"role_id": "role1", "channel_id": "chan1", "guild_id": "guild1",
			"message": "", "dev": "false", "dev_title": "true",
		}
		params, err := extractCommandParamsFunc(metaData)
		assert.NoError(t, err)
		assert.Equal(t, "role1", params.RoleID)
		assert.False(t, params.Dev)
		assert.True(t, params.DevTitle)
	})

	t.Run("Optional parameters missing", func(t *testing.T) {
		metaData := map[string]string{
			"role_id": "role1", "channel_id": "chan1", "guild_id": "guild1",
		}
		params, err := extractCommandParamsFunc(metaData)
		assert.NoError(t, err)
		assert.Equal(t, "", params.Message)
		assert.False(t, params.Dev, "Dev should default to false")
		assert.False(t, params.DevTitle, "DevTitle should default to false")
	})

	t.Run("Invalid dev parameter", func(t *testing.T) {
		metaData := map[string]string{
			"role_id": "role1", "channel_id": "chan1", "guild_id": "guild1", "dev": "not-a-bool",
		}
		params, err := extractCommandParamsFunc(metaData)
		assert.NoError(t, err)
		assert.False(t, params.Dev, "Dev should default to false on parse error")
	})

	t.Run("Invalid dev_title parameter", func(t *testing.T) {
		metaData := map[string]string{
			"role_id": "role1", "channel_id": "chan1", "guild_id": "guild1", "dev_title": "not-a-bool",
		}
		params, err := extractCommandParamsFunc(metaData)
		assert.NoError(t, err)
		assert.False(t, params.DevTitle, "DevTitle should default to false on parse error")
	})

	t.Run("Missing required parameter role_id", func(t *testing.T) {
		metaData := map[string]string{"channel_id": "chan1", "guild_id": "guild1"}
		_, err := extractCommandParamsFunc(metaData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to extract command params")
	})

	t.Run("Missing required parameter channel_id", func(t *testing.T) {
		metaData := map[string]string{"role_id": "role1", "guild_id": "guild1"}
		_, err := extractCommandParamsFunc(metaData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to extract command params")
	})

	t.Run("Missing required parameter guild_id", func(t *testing.T) {
		metaData := map[string]string{"role_id": "role1", "channel_id": "chan1"}
		_, err := extractCommandParamsFunc(metaData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to extract command params")
	})
}

func TestFetchMembersWithRoleLogic(t *testing.T) {
	guildID := "g1"
	roleID := "r1"
	channelID := "c1"
	member1 := &discordgo.Member{User: &discordgo.User{ID: "u1"}, Roles: []string{roleID}}
	member2 := &discordgo.Member{User: &discordgo.User{ID: "u2"}, Roles: []string{"other"}}
	membersListInput := []*discordgo.Member{member1, member2}
	expectedOutput := []*discordgo.Member{member1}
	mockErr := errors.New("fetch error")
	var emptyMemberList []*discordgo.Member

	t.Run("Success - Members found", func(t *testing.T) {
		mockSession := new(MockDiscordSession)

		mockSession.On("GuildMembers", guildID, "", 1000).Return(membersListInput, nil).Once()

		members, err := fetchMembersWithRoleFunc(mockSession, guildID, roleID, channelID)

		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, members)
		mockSession.AssertExpectations(t)
		mockSession.AssertNotCalled(t, "ChannelMessageSend", mock.Anything, mock.Anything)
	})

	t.Run("Success - No members with role found", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		mockSession.On("GuildMembers", guildID, "", 1000).Return([]*discordgo.Member{member2}, nil).Once()

		members, err := fetchMembersWithRoleFunc(mockSession, guildID, roleID, channelID)

		assert.NoError(t, err)
		assert.Equal(t, emptyMemberList, members)
		mockSession.AssertExpectations(t)
		mockSession.AssertNotCalled(t, "ChannelMessageSend", mock.Anything, mock.Anything)
	})

	t.Run("Error - GuildMembers fails, sending error message succeeds", func(t *testing.T) {
		mockSession := new(MockDiscordSession)

		mockSession.On("GuildMembers", guildID, "", 1000).Return(nil, mockErr).Once()

		expectedErrorMsg := fmt.Sprintf("Failed to fetch members with role: <@&%s>. Error: %v", roleID, fmt.Errorf("failed to fetch guild members: %w", mockErr)) // Match the wrapped error msg
		mockSession.On("ChannelMessageSend", channelID, expectedErrorMsg).Return(&discordgo.Message{}, nil).Once()                                                // Simulate send success

		members, err := fetchMembersWithRoleFunc(mockSession, guildID, roleID, channelID)

		assert.Error(t, err)
		assert.ErrorContains(t, err, mockErr.Error())
		assert.ErrorContains(t, err, "failed to fetch guild members")
		assert.Nil(t, members)
		mockSession.AssertExpectations(t)
	})

	t.Run("Error - GuildMembers fails, sending error message fails", func(t *testing.T) {
		mockSession := new(MockDiscordSession)
		sendErr := errors.New("send failed")
		mockSession.On("GuildMembers", guildID, "", 1000).Return(nil, mockErr).Once() // Simulate fetch failure

		expectedErrorMsg := fmt.Sprintf("Failed to fetch members with role: <@&%s>. Error: %v", roleID, fmt.Errorf("failed to fetch guild members: %w", mockErr)) // Match wrapped error
		mockSession.On("ChannelMessageSend", channelID, expectedErrorMsg).Return(nil, sendErr).Once()                                                             // Simulate send failure

		members, err := fetchMembersWithRoleFunc(mockSession, guildID, roleID, channelID)

		assert.Error(t, err)
		assert.ErrorContains(t, err, mockErr.Error())
		assert.ErrorContains(t, err, "failed to fetch guild members")
		assert.Nil(t, members)
		mockSession.AssertExpectations(t)
	})
}

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

func TestHandleDevTitleModeLogic(t *testing.T) {
	mockErr := errors.New("send failed")
	params := CommandParams{ChannelID: "c1", RoleID: "r1"}
	mentions := []string{"<@u1>", "<@u2>"}
	expectedResponse := utils.FormatDevTitleResponse(mentions, params.RoleID)

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
