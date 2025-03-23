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

}
