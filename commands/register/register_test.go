package register

import (
	constants "github.com/Real-Dev-Squad/discord-service/commands"
	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/utils"
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
	commandsCreated          []*discordgo.ApplicationCommand
	closeCalled              bool
	getUserIdCalled          bool
}

func (m *mockSession) Open() error {
	return m.openError
}

func (m *mockSession) Close() error {
	m.closeCalled = true
	return nil
}

func (m *mockSession) ApplicationCommandCreate(applicationID, guildID string, command *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error) {
	m.commandsCreated = append(m.commandsCreated, command)
	if m.commandError == nil {
		m.applicationCommandCalled = true
	}
	return nil, m.commandError
}

func (m *mockSession) GetUerId() string {
	m.getUserIdCalled = true
	return ""
}

func TestRegisterCommands(t *testing.T) {
	originalFeatureFlagState := config.AppConfig.MENTION_EACH_ENABLED
	t.Cleanup(func() {
		config.AppConfig.MENTION_EACH_ENABLED = originalFeatureFlagState
	})

	createRegisterMockSession := func(openErr, cmdErr error) *mockSession {
		sess := &mockSession{
			openError:       openErr,
			commandError:    cmdErr,
			commandsCreated: make([]*discordgo.ApplicationCommand, 0),
		}
		return sess
	}
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
	t.Run("when mention-each enabled, should attempt to register mention-each command", func(t *testing.T) {
		config.AppConfig.MENTION_EACH_ENABLED = true
		mockSess := createRegisterMockSession(nil, nil)
		assert.NotPanics(t, func() { RegisterCommands(mockSess) })

		found := false
		for _, cmd := range mockSess.commandsCreated {
			if cmd != nil && cmd.Name == utils.CommandNames.MentionEach {
				found = true
				break
			}
		}
		assert.True(t, found, "mention-each command should have been attempted for registration when flag is true")
		assert.True(t, mockSess.closeCalled)
	})

	t.Run("when mention-each disabled, should NOT attempt to register mention-each command", func(t *testing.T) {
		config.AppConfig.MENTION_EACH_ENABLED = false
		mockSess := createRegisterMockSession(nil, nil)
		assert.NotPanics(t, func() { RegisterCommands(mockSess) })

		found := false
		for _, cmd := range mockSess.commandsCreated {
			if cmd != nil && cmd.Name == utils.CommandNames.MentionEach {
				found = true
				break
			}
		}
		assert.False(t, found, "mention-each command should NOT have been attempted for registration when flag is false")
		if len(constants.Commands) > 1 {
			otherCommandExists := false
			for _, cmd := range constants.Commands {
				if cmd.Name != utils.CommandNames.MentionEach {
					otherCommandExists = true
					break
				}
			}
			if otherCommandExists {
				assert.NotEmpty(t, mockSess.commandsCreated, "Other commands should still be registered")
			} else {
				assert.Empty(t, mockSess.commandsCreated, "No commands should be registered if only mention-each exists and is disabled")
			}
		} else {
			assert.Empty(t, mockSess.commandsCreated, "No commands should be registered if only mention-each exists and is disabled")
		}
		assert.True(t, mockSess.closeCalled)
	})
}
