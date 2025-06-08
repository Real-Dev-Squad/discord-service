package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (CS *CommandHandler) verify() error {
	metaData := CS.discordMessage.MetaData
	token := metaData["token"]
	applicationId := metaData["applicationId"]

	session, err := CreateSession()
	if err != nil {
		return err
	}

	// Todo: This is a temporary solution to check that message is being sent correctly.
	content := fmt.Sprintf("Hey there <@%s>! Congratulations, you just executed your first slash command", CS.discordMessage.UserID)

	webhookEdit := &discordgo.WebhookEdit{
		Content: &content,
	}

	if _, err := session.WebhookMessageEdit(applicationId, token, "@original", webhookEdit); err != nil {
		logrus.Errorf("Error editing webhook message for application %s: %v", applicationId, err)
		return err
	}

	if err := session.Close(); err != nil {
		logrus.Errorf("Error closing session: %v", err)
		return err
	}
	return nil
}
