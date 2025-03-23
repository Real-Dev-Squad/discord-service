package handlers

import (
	"errors"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/models"
	"github.com/Real-Dev-Squad/discord-service/tests"
	_ "github.com/Real-Dev-Squad/discord-service/tests/setup"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/stretchr/testify/assert"
)

func TestMainHandler(t *testing.T) {
	t.Run("should return listeningHandler for 'listening' command", func(t *testing.T) {
		dataPacket := &dtos.DataPacket{
			CommandName: utils.CommandNames.Listening,
		}
		data, err := utils.ToByte(dataPacket)
		assert.NoError(t, err)

		handler := MainHandler(data)
		assert.NotNil(t, handler)
	})
	t.Run("should handle error and return nil for 'listening' command", func(t *testing.T) {
		dataPacket := &dtos.DataPacket{
			CommandName: utils.CommandNames.Listening,
		}
		originalFunc := utils.FromByte
		defer func() { utils.FromByte = originalFunc }()
		utils.FromByte = func(data []byte, v interface{}) error {
			return errors.New("testing error")
		}
		data, err := utils.ToByte(dataPacket)
		assert.NoError(t, err)
		handler := MainHandler(data)
		assert.Nil(t, handler)
	})
	t.Run("should return nil for invalid data", func(t *testing.T) {
		invalidData := []byte(`{"invalid": "data"}`)
		handler := MainHandler(invalidData)
		assert.Nil(t, handler)
	})
}

func TestUpdateNickname(t *testing.T) {
	t.Run("should return error if nickname is more than 32 characters", func(t *testing.T) {
		sessionWrapper := models.SessionWrapper{}
		err := UpdateNickName("userID", "nicknameWithMoreThan32Characters.", &sessionWrapper)
		assert.Error(t, err)
	})
	t.Run("should return error if nickname is more than 32 characters", func(t *testing.T) {
		sessionWrapper := models.SessionWrapper{}
		err := UpdateNickName("userID", "nicknameWithMoreThan32Characters.", &sessionWrapper)
		assert.Error(t, err)
		assert.Equal(t, "Must be 32 or fewer in length.", err.Error())
	})
	t.Run("should return error if GuildMemberNickname fails", func(t *testing.T) {
		mockSess := &tests.MockSession{GuildMemberNicknameError: true}
		err := UpdateNickName("userID", "nickname", mockSess)
		assert.Error(t, err)
	})
	t.Run("should not return error if GuildMemberNickname succeeds", func(t *testing.T) {
		mockSess := &tests.MockSession{GuildMemberNicknameError: false}
		err := UpdateNickName("userID", "nickname", mockSess)
		assert.NoError(t, err)
	})
}
