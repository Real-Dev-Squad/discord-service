package models

import (
	"testing"

	_ "github.com/Real-Dev-Squad/discord-service/tests/setup"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

var command = &discordgo.ApplicationCommand{
	Name:        utils.CommandNames.Hello,
	Description: "Greets back with hello!",
}

func TestSessionWrapper(t *testing.T) {
	mockSession := &discordgo.Session{}
	sessionWrapper := &SessionWrapper{Session: mockSession}

	t.Run("SessionWrapper should always implement open() method", func(t *testing.T) {
		assert.Panics(t, func() {
			sessionWrapper.Open()
		}, "should panic when open() is called")
	})

	t.Run("SessionWrapper should always implement close() method", func(t *testing.T) {
		assert.NotPanics(t, func() {
			sessionWrapper.Close()
		}, "should not panic when close() is called")
	})

	t.Run("SessionWrapper should always implement applicationCommandCreate() method", func(t *testing.T) {
		assert.Panics(t, func() {
			sessionWrapper.ApplicationCommandCreate("1", "2", command)
		}, "should panic when applicationCommandCreate() is called")
	})

	t.Run("SessionWrapper should always implement getUerId() method", func(t *testing.T) {
		assert.Panics(t, func() {
			sessionWrapper.GetUerId()
		}, "should panic when getUerId() is called")
	})

	t.Run("SessionWrapper should always implement guildMemberNickname() method", func(t *testing.T) {
		assert.Panics(t, func() {
			sessionWrapper.GuildMemberNickname("1", "2")
		}, "should panic when guildMemberNickname() is called")
	})

	t.Run("SessionWrapper should always implement addHandler() method", func(t *testing.T) {
		assert.NotPanics(t, func() {
			sessionWrapper.AddHandler()
		}, "should not panic when addHandler() is called")

	})

}

func TestCreateSession(t *testing.T) {
	t.Run("should create a new session", func(t *testing.T) {
		session, err := CreateSession()
		assert.NoError(t, err)
		assert.NotNil(t, session)
	})
	t.Run("should return error if creating a new session fails", func(t *testing.T) {
		originalNewDiscordSession := utils.NewDiscordSession
		defer func() { utils.NewDiscordSession = originalNewDiscordSession }()
		utils.NewDiscordSession = func(token string) (*discordgo.Session, error) {
			return nil, assert.AnError
		}
		session, err := CreateSession()
		assert.Error(t, err)
		assert.Nil(t, session)
	})
}
