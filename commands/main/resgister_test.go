package main

import (
	"testing"

	constants "github.com/Real-Dev-Squad/discord-service/commands"
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
		assert.Panics(t, main)
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
		assert.Panics(t, main)
	})
}

type mockSession struct {
	openError                error
	commandError             error
	applicationCommandCalled bool
	closeCalled              bool
	getUserIdCalled          bool
}

func (m *mockSession) open() error {
	return m.openError
}

func (m *mockSession) close() error {
	m.closeCalled = true
	return nil
}

func (m *mockSession) applicationCommandCreate(applicationID, guildID string, command *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error) {
	if m.commandError == nil {
		m.applicationCommandCalled = true
	}
	return nil, m.commandError
}

func (m *mockSession) getUerId() string {
	m.getUserIdCalled = true
	return ""
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

func TestSessionWrapper(t *testing.T) {
	mockSession := &discordgo.Session{}
	sessionWrapper := &SessionWrapper{session: mockSession}

	t.Run("SessionWrapper should always implement open() method", func(t *testing.T) {
		assert.Panics(t, func() {
			sessionWrapper.open()
		}, "should panic when open() is called")
	})

	t.Run("SessionWrapper should always implement close() method", func(t *testing.T) {
		assert.NotPanics(t, func() {
			sessionWrapper.close()
		}, "should not panic when close() is called")
	})

	t.Run("SessionWrapper should always implement applicationCommandCreate() method", func(t *testing.T) {
		assert.Panics(t, func() {
			sessionWrapper.applicationCommandCreate("1", "2", constants.Commands[0])
		}, "should panic when applicationCommandCreate() is called")
	})

	t.Run("SessionWrapper should always implement getUerId() method", func(t *testing.T) {
		assert.Panics(t, func() {
			sessionWrapper.getUerId()
		}, "should panic when getUerId() is called")
	})

}
