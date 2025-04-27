package register

import (
	"testing"

	_ "github.com/Real-Dev-Squad/discord-service/tests/helpers"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {

	t.Run("should panic when SetupConnection returns an error", func(t *testing.T) {
		originalNewDiscord := NewDiscord
		defer func() { NewDiscord = originalNewDiscord }()
		NewDiscord = func(token string) (s *discordgo.Session, err error) {
			return nil, assert.AnError
		}
		assert.Panics(t, SetupRegister)
	})
	t.Run("should call AddHandler method of session if SetupConnection succeeds", func(t *testing.T) {
		originalNewDiscord := NewDiscord
		defer func() { NewDiscord = originalNewDiscord }()
		NewDiscord = func(token string) (s *discordgo.Session, err error) {
			mockSession := &discordgo.Session{
				State: &discordgo.State{},
			}
			return mockSession, nil
		}
		assert.Panics(t, SetupRegister)
	})
}

type mockSession struct {
	openError                error
	commandError             error
	applicationCommandCalled bool
	closeCalled              bool
	getUserIdCalled          bool
	channelMessageSendCalled bool
	channelMessageSendError  error
	guildMembersCalled       bool
	guildMembersError        error
}

func (m *mockSession) Open() error {
	return m.openError
}

func (m *mockSession) Close() error {
	m.closeCalled = true
	return nil
}

func (m *mockSession) ApplicationCommandCreate(applicationID, guildID string, command *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error) {
	if m.commandError == nil {
		m.applicationCommandCalled = true
	}
	return nil, m.commandError
}

func (m *mockSession) GetUserId() string {
	m.getUserIdCalled = true
	return ""
}

func (m *mockSession) GuildMembers(guildID, after string, limit int) ([]*discordgo.Member, error) {
	m.guildMembersCalled = true
	return nil, m.guildMembersError
}
func (m *mockSession) ChannelMessageSend(channelID string, content string) (*discordgo.Message, error) {
	m.channelMessageSendCalled = true
	return nil, m.channelMessageSendError
}

func TestRegisterCommands(t *testing.T) {
	t.Run("should not panic when Open() returns no error", func(t *testing.T) {
		mockSess := &mockSession{openError: nil, commandError: nil}
		assert.NotPanics(t, func() {
			RegisterCommands(mockSess)
		}, "RegisterCommands should not panic when Open is successful")

	})
	t.Run("should panic when Open() returns an error", func(t *testing.T) {
		mockSess := &mockSession{openError: assert.AnError, commandError: nil}
		assert.Panics(t, func() {
			RegisterCommands(mockSess)
		}, "RegisterCommands should panic when Open returns an error")
	})
	t.Run("should panic when openSession.ApplicationCommandCreate() returns an error", func(t *testing.T) {
		mockSess := &mockSession{openError: nil, commandError: assert.AnError}

		assert.Panics(t, func() {
			RegisterCommands(mockSess)
		}, "RegisterCommands should panic when ApplicationCommandCreate returns an error")
	})
	t.Run("should call all methods when none of the methods returns no error", func(t *testing.T) {
		mockSess := &mockSession{openError: nil, commandError: nil}
		assert.NotPanics(t, func() {
			RegisterCommands(mockSess)
		})
		assert.True(t, mockSess.applicationCommandCalled)
		assert.True(t, mockSess.getUserIdCalled)
		assert.True(t, mockSess.closeCalled)
	})
}
