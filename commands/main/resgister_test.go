package main

import (
	"testing"

	"github.com/Real-Dev-Squad/discord-service/tests"
	_ "github.com/Real-Dev-Squad/discord-service/tests/setup"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {

	t.Run("should panic when SetupConnection returns an error", func(t *testing.T) {
		originalNewDiscord := utils.NewDiscordSession
		defer func() { utils.NewDiscordSession = originalNewDiscord }()
		utils.NewDiscordSession = func(token string) (s *discordgo.Session, err error) {
			return nil, assert.AnError
		}
		assert.Panics(t, main)
	})
	t.Run("should call AddHandler method of session if SetupConnection succeeds", func(t *testing.T) {
		originalNewDiscord := utils.NewDiscordSession
		defer func() { utils.NewDiscordSession = originalNewDiscord }()
		utils.NewDiscordSession = func(token string) (s *discordgo.Session, err error) {
			mockSession := &discordgo.Session{
				State: &discordgo.State{},
			}
			return mockSession, nil
		}
		assert.Panics(t, main)
	})
}

func TestRegisterCommands(t *testing.T) {
	t.Run("should not panic when Open() returns no error", func(t *testing.T) {
		mockSess := &tests.MockSession{OpenError: nil, CommandError: nil}
		assert.NotPanics(t, func() {
			RegisterCommands(mockSess)
		}, "RegisterCommands should not panic when Open is successful")

	})
	t.Run("should panic when Open() returns an error", func(t *testing.T) {
		mockSess := &tests.MockSession{OpenError: assert.AnError, CommandError: nil}
		assert.Panics(t, func() {
			RegisterCommands(mockSess)
		}, "RegisterCommands should panic when Open returns an error")
	})
	t.Run("should panic when openSession.ApplicationCommandCreate() returns an error", func(t *testing.T) {
		mockSess := &tests.MockSession{OpenError: nil, CommandError: assert.AnError}

		assert.Panics(t, func() {
			RegisterCommands(mockSess)
		}, "RegisterCommands should panic when ApplicationCommandCreate returns an error")
	})
	t.Run("should panic when openSession.ApplicationCommandCreate() returns an error", func(t *testing.T) {
		mockSess := &tests.MockSession{OpenError: nil, CommandError: assert.AnError}

		assert.Panics(t, func() {
			RegisterCommands(mockSess)
		}, "RegisterCommands should panic when ApplicationCommandCreate returns an error")
	})
	t.Run("should call all methods when none of the methods returns no error", func(t *testing.T) {
		mockSess := &tests.MockSession{OpenError: nil, CommandError: nil}
		assert.NotPanics(t, func() {
			RegisterCommands(mockSess)
		})
		assert.True(t, mockSess.ApplicationCommandCalled)
		assert.True(t, mockSess.GetUserIdCalled)
		assert.True(t, mockSess.CloseCalled)
	})
}
