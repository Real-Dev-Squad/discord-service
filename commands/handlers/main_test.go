package handlers

import (
	"errors"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	_ "github.com/Real-Dev-Squad/discord-service/tests/helpers"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestMainHandler(t *testing.T) {
	t.Run("should return listeningHandler for 'listening' command", func(t *testing.T) {
		dataPacket := &dtos.DataPacket{
			CommandName: utils.CommandNames.Listening,
		}
		data, err := dtos.ToByte(dataPacket)
		assert.NoError(t, err)

		handler := MainHandler(data)
		assert.NotNil(t, handler)
	})
	t.Run("should return nil for invalid json data", func(t *testing.T) {
		invalidData := []byte(`{"userId":"1234567890","cmdName":"listening"}`)
		handler := MainHandler(invalidData)
		assert.Nil(t, handler)
	})

	t.Run("Should return nil for unknown commands", func(t *testing.T) {
		dp := &dtos.DataPacket{
			CommandName: "unknown",
		}
		data, err := dtos.ToByte(dp)
		assert.NoError(t, err)

		handler := MainHandler(data)
		assert.Nil(t, handler)
	})
}

func TestCreateSession(t *testing.T) {
	t.Run("should fail if NewDiscord returns an error", func(t *testing.T) {
		originalNewDiscord := NewDiscord
		defer func() { NewDiscord = originalNewDiscord }()
		NewDiscord = func(token string) (s *discordgo.Session, err error) {
			return nil, errors.New("testing error")
		}
		_, err := CreateSession()
		assert.Error(t, err)
	})
	t.Run("should initiate open session if NewDiscord returns no error", func(t *testing.T) {
		originalNewDiscord := NewDiscord
		defer func() { NewDiscord = originalNewDiscord }()
		NewDiscord = func(token string) (s *discordgo.Session, err error) {
			return &discordgo.Session{}, nil
		}
		assert.Panics(t, func() { CreateSession() })
	})
}

func TestUpdateNickName(t *testing.T) {

	var originalCreateSession = CreateSession
	defer func() { CreateSession = originalCreateSession }()

	t.Run("should return error if newNickName is longer than 32 characters", func(t *testing.T) {
		err := UpdateNickName("userID", "ThisIsAVeryLongNicknameThatExceedsTheLimit")
		assert.Error(t, err)
		assert.Equal(t, "Must be 32 or fewer in length.", err.Error())
	})

	t.Run("should return error if CreateSession fails", func(t *testing.T) {
		CreateSession = func() (DiscordSessionWrapper, error) {
			return nil, errors.New("failed to create session")
		}
		err := UpdateNickName("userID", "validNickname")
		assert.Error(t, err)
		assert.Equal(t, "failed to create session", err.Error())
	})
	t.Run("should hit GuildMemberNickname if CreateSession succeeds", func(t *testing.T) {
		CreateSession = func() (DiscordSessionWrapper, error) {
			panic("GuildMemberNickname called")
		}
		assert.Panics(t, func() { UpdateNickName("userID", "validNickname") })

	})

}
