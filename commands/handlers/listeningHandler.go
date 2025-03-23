package handlers

import (
	"fmt"
	"strings"

	"github.com/Real-Dev-Squad/discord-service/models"
	"github.com/Real-Dev-Squad/discord-service/utils"
)

func (s *CommandHandler) listeningHandler() error {
	metaData := s.discordMessage.MetaData
	nickName := metaData["nickname"]
	if metaData["value"] == "true" {
		nickName = fmt.Sprintf("%s%s%s", utils.NICKNAME_PREFIX, nickName, utils.NICKNAME_SUFFIX)
	} else {
		nickName = strings.TrimPrefix(strings.TrimSuffix(nickName, utils.NICKNAME_SUFFIX), utils.NICKNAME_PREFIX)
	}
	sessionWrapper, err := models.CreateSession()
	if err != nil {
		return err
	}
	return UpdateNickName(s.discordMessage.UserID, nickName, sessionWrapper)
}
