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
		data, err := dataPacket.ToByte()
		assert.NoError(t, err)

		handler := MainHandler(data)
		assert.NotNil(t, handler)
	})
	t.Run("should return mentionEachHandler for 'mention-each' command", func(t *testing.T) {
		dataPacket := &dtos.DataPacket{
			CommandName: utils.CommandNames.MentionEach,
		}

		data, err := dataPacket.ToByte()
		assert.NoError(t, err)

		handler := MainHandler(data)
		assert.NotNil(t, handler, "Handler should not be nil for mention-each command")

	})
	t.Run("should return nil for invalid data", func(t *testing.T) {
		invalidData := []byte(`{"invalid": "data"}`)
		handler := MainHandler(invalidData)
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

func mockCreateSession() (*discordgo.Session, error) {
	return &discordgo.Session{}, nil
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
		CreateSession = func() (*discordgo.Session, error) {
			return nil, errors.New("failed to create session")
		}
		err := UpdateNickName("userID", "validNickname")
		assert.Error(t, err)
		assert.Equal(t, "failed to create session", err.Error())
	})
	t.Run("should hit GuildMemberNickname if CreateSession succeeds", func(t *testing.T) {
		CreateSession = mockCreateSession
		assert.Panics(t, func() { UpdateNickName("userID", "validNickname") })
	})

}
