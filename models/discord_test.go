package models

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

const helloCommandName = "hello"

var command = &discordgo.ApplicationCommand{
	Name:        helloCommandName,
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
			sessionWrapper.GetUserId()
		}, "should panic when getUerId() is called")
	})
	t.Run("SessionWrapper should implement GuildMembers() method", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = sessionWrapper.GuildMembers("guild1", "", 100)
		}, "should panic when GuildMembers() is called on dummy session")
	})

	t.Run("SessionWrapper should implement ChannelMessageSend() method", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = sessionWrapper.ChannelMessageSend("chan1", "hello")
		}, "should panic when ChannelMessageSend() is called on dummy session")
	})

}
