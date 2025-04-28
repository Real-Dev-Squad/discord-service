package handlers

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/models"
	"github.com/Real-Dev-Squad/discord-service/tests/mocks"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

	member1 := discordgo.Member{User: &discordgo.User{ID: "user1"}, Roles: []string{"testRole"}}
	membersList := []discordgo.Member{member1}
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
		fetchMembersWithRoleFunc = func(session models.SessionInterface, guildID, roleID, channelID string) ([]discordgo.Member, error) {
			assert.NotNil(t, session)
			assert.Equal(t, "testGuild", guildID)
			return membersList, nil
		}
		handleStandardModeFunc = func(session models.SessionInterface, mentions []string, params CommandParams) error {
			assert.NotNil(t, session)
			assert.Contains(t, mentions, "<@user1>")
			assert.False(t, params.Dev)
			assert.False(t, params.DevTitle)
			return nil
		}
		handleDevModeFunc = func(models.SessionInterface, []string, CommandParams) error {
			t.Fatal("handleDevMode called unexpectedly")
			return nil
		}
		handleDevTitleModeFunc = func(models.SessionInterface, []string, CommandParams) error {
			t.Fatal("handleDevTitleMode called unexpectedly")
			return nil
		}
		sendNoMembersMessageFunc = func(models.SessionInterface, string) error {
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
				RoleID: "testRole", ChannelID: "testChannel", GuildID: "testGuild", Message: "Dev Msg", Dev: true, DevTitle: false,
			}, nil
		}
		CreateSession = func() (*discordgo.Session, error) { return &discordgo.Session{}, nil }
		fetchMembersWithRoleFunc = func(models.SessionInterface, string, string, string) ([]discordgo.Member, error) {
			return membersList, nil
		}
		handleDevModeFunc = func(session models.SessionInterface, mentions []string, params CommandParams) error {
			assert.NotNil(t, session)
			assert.True(t, params.Dev)
			assert.Equal(t, "Dev Msg", params.Message)
			assert.Contains(t, mentions, "<@user1>")
			return nil
		}

		handleStandardModeFunc = func(models.SessionInterface, []string, CommandParams) error {
			t.Fatal("handleStandardMode called unexpectedly")
			return nil
		}
		handleDevTitleModeFunc = func(models.SessionInterface, []string, CommandParams) error {
			t.Fatal("handleDevTitleMode called unexpectedly")
			return nil
		}
		sendNoMembersMessageFunc = func(models.SessionInterface, string) error {
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
				RoleID: "testRole", ChannelID: "testChannel", GuildID: "testGuild", Message: "", Dev: false, DevTitle: true,
			}, nil
		}
		CreateSession = func() (*discordgo.Session, error) { return &discordgo.Session{}, nil }
		fetchMembersWithRoleFunc = func(models.SessionInterface, string, string, string) ([]discordgo.Member, error) {
			return membersList, nil
		}
		handleDevTitleModeFunc = func(session models.SessionInterface, mentions []string, params CommandParams) error {
			assert.NotNil(t, session)
			assert.True(t, params.DevTitle)
			assert.Contains(t, mentions, "<@user1>")
			return nil
		}
		handleStandardModeFunc = func(models.SessionInterface, []string, CommandParams) error {
			t.Fatal("handleStandardMode called unexpectedly")
			return nil
		}
		handleDevModeFunc = func(models.SessionInterface, []string, CommandParams) error {
			t.Fatal("handleDevMode called unexpectedly")
			return nil
		}
		sendNoMembersMessageFunc = func(models.SessionInterface, string) error {
			t.Fatal("sendNoMembersMessage called unexpectedly")
			return nil
		}

		commandHandler := setupTestCommandHandler(nil)
		err := commandHandler.mentionEachHandler()
		assert.NoError(t, err)
	})

	t.Run("Success No Members Found", func(t *testing.T) {

		extractCommandParamsFunc = func(map[string]string) (CommandParams, error) {
			return CommandParams{ChannelID: "testChannel"}, nil
		}

		CreateSession = func() (*discordgo.Session, error) { return &discordgo.Session{}, nil }
		fetchMembersWithRoleFunc = func(models.SessionInterface, string, string, string) ([]discordgo.Member, error) {
			emptyMembersList := []discordgo.Member{}
			return emptyMembersList, nil
		}

		sendNoMembersMessageFunc = func(session models.SessionInterface, channelID string) error {
			assert.NotNil(t, session)
			assert.Equal(t, "testChannel", channelID)
			return nil
		}

		handleStandardModeFunc = func(models.SessionInterface, []string, CommandParams) error {
			t.Fatal("handleStandardMode called unexpectedly")
			return nil
		}

		handleDevModeFunc = func(models.SessionInterface, []string, CommandParams) error {
			t.Fatal("handleDevMode called unexpectedly")
			return nil
		}

		handleDevTitleModeFunc = func(models.SessionInterface, []string, CommandParams) error {
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
		fetchMembersWithRoleFunc = func(models.SessionInterface, string, string, string) ([]discordgo.Member, error) {
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
		fetchMembersWithRoleFunc = func(models.SessionInterface, string, string, string) ([]discordgo.Member, error) {
			emptyMembersList := []discordgo.Member{}
			return emptyMembersList, nil
		}
		sendNoMembersMessageFunc = func(models.SessionInterface, string) error {
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
		fetchMembersWithRoleFunc = func(models.SessionInterface, string, string, string) ([]discordgo.Member, error) {
			emptyMembersList := []discordgo.Member{}
			return emptyMembersList, nil
		}
		handleStandardModeFunc = func(models.SessionInterface, []string, CommandParams) error {
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
		fetchMembersWithRoleFunc = func(models.SessionInterface, string, string, string) ([]discordgo.Member, error) {
			membersList := []discordgo.Member{}
			return membersList, nil
		}
		handleDevModeFunc = func(models.SessionInterface, []string, CommandParams) error {
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
		fetchMembersWithRoleFunc = func(models.SessionInterface, string, string, string) ([]discordgo.Member, error) {
			membersList := []discordgo.Member{}
			return membersList, nil
		}
		handleDevTitleModeFunc = func(models.SessionInterface, []string, CommandParams) error {
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

		mockSessionInstance := new(mocks.DiscordSession)
		closeErr := errors.New("failed t0 close session")
		mockSessionInstance.On("Close").Return(closeErr).Once()

		CreateSession = func() (*discordgo.Session, error) {
			return &discordgo.Session{}, nil
		}

		fetchMembersWithRoleFunc = func(session models.SessionInterface, guildID, roleID, channelID string) ([]discordgo.Member, error) {
			return membersList, nil
		}

		handleStandardModeFunc = func(session models.SessionInterface, mentions []string, params CommandParams) error {
			return nil
		}

		handleDevModeFunc = func(session models.SessionInterface, mentions []string, params CommandParams) error {
			return nil
		}
		handleDevTitleModeFunc = func(models.SessionInterface, []string, CommandParams) error {
			t.Fatal("handleDevTitleMode called unexpectedly")
			return nil
		}
		sendNoMembersMessageFunc = func(models.SessionInterface, string) error {
			t.Fatal("sendNoMembersMessage called unexpectedly")
			return nil
		}

		commandHandler := setupTestCommandHandler(nil)
		err := commandHandler.mentionEachHandler()
		assert.NoError(t, err)
	})
}

func TestExtractCommandParamsLogic(t *testing.T) {

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
		}
		params, err := extractCommandParamsFunc(metaData)
		assert.NoError(t, err)
		assert.Equal(t, "role1", params.RoleID)
		assert.Equal(t, "", params.Message)
		assert.False(t, params.Dev)
		assert.False(t, params.DevTitle)
	})

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
			"role_id": "role1", "channel_id": "chan1", "guild_id": "guild1", "dev_title": "123",
		}
		params, err := extractCommandParamsFunc(metaData)
		assert.NoError(t, err)
		assert.False(t, params.DevTitle, "DevTitle should default to false on parse error")
	})

	expectedErrorMsg := "failed to extract command params: missing or empty role_id, channel_id, or guild_id"

	t.Run("Missing required parameter role_id", func(t *testing.T) {
		metaData := map[string]string{"channel_id": "chan1", "guild_id": "guild1"}
		_, err := extractCommandParamsFunc(metaData)
		assert.Error(t, err)
		assert.EqualError(t, err, expectedErrorMsg)
	})

	t.Run("Empty required parameter role_id", func(t *testing.T) {
		metaData := map[string]string{"role_id": "", "channel_id": "chan1", "guild_id": "guild1"}
		_, err := extractCommandParamsFunc(metaData)
		assert.Error(t, err)
		assert.EqualError(t, err, expectedErrorMsg)
	})

	t.Run("Missing required parameter channel_id", func(t *testing.T) {
		metaData := map[string]string{"role_id": "role1", "guild_id": "guild1"}
		_, err := extractCommandParamsFunc(metaData)
		assert.Error(t, err)
		assert.EqualError(t, err, expectedErrorMsg)
	})

	t.Run("Empty required parameter channel_id", func(t *testing.T) {
		metaData := map[string]string{"role_id": "role1", "channel_id": "", "guild_id": "guild1"}
		_, err := extractCommandParamsFunc(metaData)
		assert.Error(t, err)
		assert.EqualError(t, err, expectedErrorMsg)
	})

	t.Run("Missing required parameter guild_id", func(t *testing.T) {
		metaData := map[string]string{"role_id": "role1", "channel_id": "chan1"}
		_, err := extractCommandParamsFunc(metaData)
		assert.Error(t, err)
		assert.EqualError(t, err, expectedErrorMsg)
	})

	t.Run("Empty required parameter guild_id", func(t *testing.T) {
		metaData := map[string]string{"role_id": "role1", "channel_id": "chan1", "guild_id": ""}
		_, err := extractCommandParamsFunc(metaData)
		assert.Error(t, err)
		assert.EqualError(t, err, expectedErrorMsg)
	})
}

func TestFetchMembersWithRoleLogic(t *testing.T) {
	guildID := "g1"
	roleID := "r1"
	channelID := "c1"
	member1 := &discordgo.Member{User: &discordgo.User{ID: "u1"}, Roles: []string{roleID}}
	member2 := &discordgo.Member{User: &discordgo.User{ID: "u2"}, Roles: []string{"other"}}
	membersListInputSuccess := []*discordgo.Member{member1, member2}
	expectedOutputSuccess := []discordgo.Member{*member1}
	mockErr := errors.New("API error from GuildMembers")
	var emptyMemberList []*discordgo.Member

	t.Run("Success - Members found", func(t *testing.T) {
		mockSession := new(mocks.DiscordSession)

		mockSession.On("GuildMembers", guildID, "", 1000).Return(membersListInputSuccess, nil).Once()
		mockSession.On("GuildMembers", guildID, member2.User.ID, 1000).Return(emptyMemberList, nil).Once()

		members, err := fetchMembersWithRoleFunc(mockSession, guildID, roleID, channelID)

		assert.NoError(t, err)
		assert.Equal(t, expectedOutputSuccess, members)
		mockSession.AssertExpectations(t)
		mockSession.AssertNotCalled(t, "ChannelMessageSend", mock.Anything, mock.Anything)
	})

	t.Run("Success - No members with role found", func(t *testing.T) {
		mockSession := new(mocks.DiscordSession)
		membersInputNoMatch := []*discordgo.Member{member2}

		mockSession.On("GuildMembers", guildID, "", 1000).Return(membersInputNoMatch, nil).Once()
		mockSession.On("GuildMembers", guildID, member2.User.ID, 1000).Return(emptyMemberList, nil).Once()

		members, err := fetchMembersWithRoleFunc(mockSession, guildID, roleID, channelID)

		assert.NoError(t, err)
		assert.Empty(t, members)
		mockSession.AssertExpectations(t)
		mockSession.AssertNotCalled(t, "ChannelMessageSend", mock.Anything, mock.Anything)
	})

	t.Run("Error - GuildMembers fails, sending error message succeeds", func(t *testing.T) {
		mockSession := new(mocks.DiscordSession)

		mockSession.On("GuildMembers", guildID, "", 1000).Return(nil, mockErr).Once()

		expectedErrorMsg := fmt.Sprintf("Sorry, I couldn't fetch members for role <@&%s> right now. Please try again later.", roleID)
		mockSession.On("ChannelMessageSend", channelID, expectedErrorMsg).Return(&discordgo.Message{}, nil).Once()

		members, err := fetchMembersWithRoleFunc(mockSession, guildID, roleID, channelID)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to fetch guild members chunk:")
		assert.ErrorContains(t, err, mockErr.Error())
		assert.Nil(t, members)
		mockSession.AssertExpectations(t)
	})

	t.Run("Error - GuildMembers fails, sending error message fails", func(t *testing.T) {
		mockSession := new(mocks.DiscordSession)
		sendErr := errors.New("send failed")

		mockSession.On("GuildMembers", guildID, "", 1000).Return(nil, mockErr).Once()

		expectedErrorMsg := fmt.Sprintf("Sorry, I couldn't fetch members for role <@&%s> right now. Please try again later.", roleID)
		mockSession.On("ChannelMessageSend", channelID, expectedErrorMsg).Return(nil, sendErr).Once()

		members, err := fetchMembersWithRoleFunc(mockSession, guildID, roleID, channelID)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to fetch guild members chunk:")
		assert.ErrorContains(t, err, mockErr.Error())
		assert.Nil(t, members)
		mockSession.AssertExpectations(t)
	})
}

func TestSendNoMembersMessageLogic(t *testing.T) {
	channelID := "c1"
	expectedMsg := "Sorry, no members found with this role"
	mockErr := errors.New("send failed")

	t.Run("Success sending message", func(t *testing.T) {
		mockSession := new(mocks.DiscordSession)
		mockSession.On("ChannelMessageSend", channelID, expectedMsg).Return(&discordgo.Message{}, nil).Once()

		err := sendNoMembersMessageFunc(mockSession, channelID)

		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Error sending message", func(t *testing.T) {
		mockSession := new(mocks.DiscordSession)
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
	mentions := []string{"<@u1>", "<@u2>", "<@u3>", "<@u4>", "<@u5>", "<@u6>"}

	t.Run("Success with batching", func(t *testing.T) {
		mockSession := new(mocks.DiscordSession)
		for i := 0; i < 6; i++ {
			expectedMsg := fmt.Sprintf("%s %s", params.Message, mentions[i])
			mockSession.On("ChannelMessageSend", params.ChannelID, expectedMsg).Return(&discordgo.Message{}, nil).Once()
		}

		err := handleDevModeFunc(mockSession, mentions, params)

		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Success single batch", func(t *testing.T) {
		mockSession := new(mocks.DiscordSession)
		shortMentions := []string{"<@uA>", "<@uB>"}
		mockSession.On("ChannelMessageSend", params.ChannelID, fmt.Sprintf("%s %s", params.Message, shortMentions[0])).Return(&discordgo.Message{}, nil).Once()
		mockSession.On("ChannelMessageSend", params.ChannelID, fmt.Sprintf("%s %s", params.Message, shortMentions[1])).Return(&discordgo.Message{}, nil).Once()

		err := handleDevModeFunc(mockSession, shortMentions, params)

		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Success no mentions", func(t *testing.T) {
		mockSession := new(mocks.DiscordSession)
		err := handleDevModeFunc(mockSession, []string{}, params)
		assert.NoError(t, err)
		mockSession.AssertNotCalled(t, "ChannelMessageSend", mock.Anything, mock.Anything)
	})

	t.Run("Success no custom message", func(t *testing.T) {
		mockSession := new(mocks.DiscordSession)
		localParams := CommandParams{ChannelID: "c1", Message: "", RoleID: "r1"}
		mention := "<@uOnly>"
		mockSession.On("ChannelMessageSend", localParams.ChannelID, mention).Return(&discordgo.Message{}, nil).Once()

		err := handleDevModeFunc(mockSession, []string{mention}, localParams)

		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Error during sending some mentions", func(t *testing.T) {
		mockSession := new(mocks.DiscordSession)
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
		mockSession := new(mocks.DiscordSession)
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedResponse).Return(&discordgo.Message{}, nil).Once()

		err := handleDevTitleModeFunc(mockSession, mentions, params)

		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Error sending message", func(t *testing.T) {
		mockSession := new(mocks.DiscordSession)
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
		mockSession := new(mocks.DiscordSession)
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedResponse).Return(&discordgo.Message{}, nil).Once()

		err := handleStandardModeFunc(mockSession, mentions, params)

		assert.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("Error sending message", func(t *testing.T) {
		mockSession := new(mocks.DiscordSession)
		mockSession.On("ChannelMessageSend", params.ChannelID, expectedResponse).Return(nil, mockErr).Once()

		err := handleStandardModeFunc(mockSession, mentions, params)

		assert.Error(t, err)
		assert.Equal(t, mockErr, err)
		mockSession.AssertExpectations(t)
	})
}
