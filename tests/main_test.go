package tests

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func Test_MockSessionHelper(t *testing.T) {

	t.Run("Open", func(t *testing.T) {
		session := &MockSession{}
		err := session.Open()
		assert.NoError(t, err)
	})

	t.Run("Close", func(t *testing.T) {
		session := &MockSession{}
		err := session.Close()
		assert.NoError(t, err)
		assert.True(t, session.CloseCalled)
	})

	t.Run("ApplicationCommandCreate", func(t *testing.T) {
		session := &MockSession{}
		_, err := session.ApplicationCommandCreate("", "", &discordgo.ApplicationCommand{})
		assert.NoError(t, err)
		assert.True(t, session.ApplicationCommandCalled)
	})

	t.Run("GetUserId", func(t *testing.T) {
		session := &MockSession{}
		userId := session.GetUserId()
		assert.Equal(t, "", userId)
		assert.True(t, session.GetUserIdCalled)
	})

	t.Run("GuildMemberNickname", func(t *testing.T) {
		session := &MockSession{}
		t.Run("should return error if GuildMemberNickname fails", func(t *testing.T) {
			session.ForceErrorForGuildMemberNickname = true
			err := session.GuildMemberNickname("", "")
			assert.Error(t, err)
		})
		t.Run("should not return error if GuildMemberNickname succeeds", func(t *testing.T) {
			session.ForceErrorForGuildMemberNickname = false
			err := session.GuildMemberNickname("", "")
			assert.NoError(t, err)
		})
	})
}
