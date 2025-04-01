package handlers

import (
	"fmt"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const (
	BatchSize  = 5
	BatchDelay = 1 * time.Second
)

type CommandParams struct {
	RoleID    string
	ChannelID string
	GuildID   string
	Message   string
	Dev       bool
	DevTitle  bool
}

type DiscordSessionWrapper struct {
	*discordgo.Session
}

func (s *DiscordSessionWrapper) GuildMembers(guildID string, after string, limit int) ([]*discordgo.Member, error) {
	members, err := s.Session.GuildMembers(guildID, after, limit)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (s *DiscordSessionWrapper) ChannelMessageSend(channelID, content string) (*discordgo.Message, error) {
	msg, err := s.Session.ChannelMessageSend(channelID, content)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

var _ utils.DiscordSessionInterface = (*DiscordSessionWrapper)(nil)

func (s *CommandHandler) mentionEachHandler() error {
	logrus.Info("Processing mention-each command")

	params, err := extractCommandParams(s.discordMessage.MetaData)
	if err != nil {
		logrus.Errorf("Parameter extraction failed: %v", err)
		return err
	}

	discordSession, err := CreateSession()
	if err != nil {
		logrus.Errorf("Error creating session: %v", err)
		return err
	}

	session := &DiscordSessionWrapper{discordSession}

	defer func() {
		if closeErr := session.Close(); closeErr != nil {
			logrus.Errorf("Error closing session: %v", closeErr)
		}
	}()

	members, err := fetchMembersWithRole(session, params.GuildID, params.RoleID, params.ChannelID)
	if err != nil {
		return err
	}

	if len(members) == 0 {
		return sendNoMembersMessage(session, params.ChannelID)
	}

	mentions := utils.FormatUserMentions(members)

	if params.DevTitle {
		return handleDevTitleMode(session, mentions, params)
	} else if params.Dev {
		return handleDevMode(session, mentions, params)
	} else {
		return handleStandardMode(session, mentions, params)
	}
}

func extractCommandParams(metaData map[string]string) (CommandParams, error) {
	params := CommandParams{
		RoleID:    metaData["role_id"],
		ChannelID: metaData["channel_id"],
		GuildID:   metaData["guild_id"],
		Message:   metaData["message"],
	}

	if params.RoleID == "" || params.ChannelID == "" || params.GuildID == "" {
		return params, fmt.Errorf("failed to extract command params")
	}

	if devStr := metaData["dev"]; devStr != "" {
		dev, err := strconv.ParseBool(devStr)
		if err == nil {
			params.Dev = dev
		}
	}

	if devTitleStr := metaData["dev_title"]; devTitleStr != "" {
		devTitle, err := strconv.ParseBool(devTitleStr)
		if err == nil {
			params.DevTitle = devTitle
		}
	}

	return params, nil
}
func fetchMembersWithRole(session utils.DiscordSessionInterface, guildID, roleID, channelID string) ([]*discordgo.Member, error) {
	members, err := utils.GetUsersWithRole(session, guildID, roleID)
	if err != nil {
		logrus.Errorf("Failed to fetch members with role: %v", err)
		errorMsg := fmt.Sprintf("Failed to fetch members with role: %v", err)
		_, sendErr := session.ChannelMessageSend(channelID, errorMsg)
		if sendErr != nil {
			logrus.Errorf("Failed to send error message: %v", sendErr)
		}
		return nil, err
	}

	return members, nil
}

func sendNoMembersMessage(session utils.DiscordSessionInterface, channelID string) error {
	_, err := session.ChannelMessageSend(channelID, "Sorry, no members found with this role")
	if err != nil {
		logrus.Errorf("Failed to send empty response: %v", err)
		return err
	}
	return nil
}
func handleDevMode(session utils.DiscordSessionInterface, mentions []string, params CommandParams) error {
	for i := 0; i < len(mentions); i += BatchSize {
		end := i + BatchSize
		if end > len(mentions) {
			end = len(mentions)
		}

		// Process current batch
		for j := i; j < end; j++ {
			msgContent := mentions[j]
			if params.Message != "" {
				msgContent = fmt.Sprintf("%s %s", params.Message, mentions[j])
			}

			_, err := session.ChannelMessageSend(params.ChannelID, msgContent)
			if err != nil {
				logrus.Errorf("Failed to send individual mention: %v", err)
				return err
			}
		}

		// Rate limiting between batches
		if end < len(mentions) {
			logrus.Infof("Rate limiting between batches")
			time.Sleep(BatchDelay)
		}
	}

	return nil
}
func handleDevTitleMode(session utils.DiscordSessionInterface, mentions []string, params CommandParams) error {
	response := utils.FormatDevTitleResponse(mentions, params.RoleID)
	_, err := session.ChannelMessageSend(params.ChannelID, response)
	if err != nil {
		logrus.Errorf("Failed to send dev_title response: %v", err)
		return err
	}
	return nil
}
func handleStandardMode(session utils.DiscordSessionInterface, mentions []string, params CommandParams) error {
	response := utils.FormatMentionResponse(mentions, params.Message)
	_, err := session.ChannelMessageSend(params.ChannelID, response)
	if err != nil {
		logrus.Errorf("Failed to send mention response: %v", err)
		return err
	}

	logrus.Infof("Successfully processed mention-each command for role: %s", params.RoleID)
	return nil
}
