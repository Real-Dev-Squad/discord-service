package handlers

import (
	"fmt"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const (
	BatchSize            = 5
	BatchDelay           = 1 * time.Second
	MaxUserMessageLength = 1000
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

var (
	extractCommandParamsFunc = func(metaData map[string]string) (CommandParams, error) {
		params := CommandParams{}
		roleID, roleOK := metaData["role_id"]
		channelID, channelOK := metaData["channel_id"]
		guildID, guildOK := metaData["guild_id"]

		if !roleOK || roleID == "" || !channelOK || channelID == "" || !guildOK || guildID == "" {
			logrus.WithFields(logrus.Fields{
				"role_id_present":    roleOK,
				"channel_id_present": channelOK,
				"guild_id_present":   guildOK,
				"role_id_value":      metaData["role_id"],
				"channel_id_value":   metaData["channel_id"],
				"guild_id_value":     metaData["guild_id"],
				"metadata":           metaData,
			}).Error("Missing or empty required parameters (role_id, channel_id, guild_id) in metadata")
			return params, fmt.Errorf("failed to extract command params: missing or empty role_id, channel_id, or guild_id")
		}
		params.RoleID = roleID
		params.ChannelID = channelID
		params.GuildID = guildID

		userMessage := metaData["message"]
		if len(userMessage) > MaxUserMessageLength {
			truncatedSuffix := "..."
			logrus.Warnf("User provided message length (%d) exceeds limit (%d). Truncating.", len(userMessage), MaxUserMessageLength)
			userMessage = userMessage[:MaxUserMessageLength] + truncatedSuffix
		}
		params.Message = userMessage

		if devStr := metaData["dev"]; devStr != "" {
			dev, err := strconv.ParseBool(devStr)
			if err == nil {
				params.Dev = dev
			} else {
				logrus.Warnf("Invalid boolean value for 'dev' flag: '%s'. Defaulting to false.", devStr)
			}
		}
		if devTitleStr := metaData["dev_title"]; devTitleStr != "" {
			devTitle, err := strconv.ParseBool(devTitleStr)
			if err == nil {
				params.DevTitle = devTitle
			} else {
				logrus.Warnf("Invalid boolean value for 'dev_title' flag: '%s'. Defaulting to false.", devTitleStr)
			}
		}

		return params, nil
	}

	fetchMembersWithRoleFunc = func(session utils.DiscordSessionInterface, guildID, roleID, channelID string) ([]*discordgo.Member, error) {
		members, err := utils.GetUsersWithRole(session, guildID, roleID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"guild_id":   guildID,
				"role_id":    roleID,
				"channel_id": channelID,
			}).Errorf("GetUserWithRole failed within fetchMembersWithRole: %v", err)
			errorMsg := fmt.Sprintf("Sorry, I couldn't fetch members for role <@&%s> right now. Please try again later.", roleID)
			_, sendErr := session.ChannelMessageSend(channelID, errorMsg)
			if sendErr != nil {
				logrus.WithFields(logrus.Fields{
					"originalError": err,
					"sendError":     sendErr,
					"channelID":     channelID,
				}).Errorf("Failed to send error message: %v", sendErr)
			}
			return nil, err
		}
		return members, nil
	}

	sendNoMembersMessageFunc = func(session utils.DiscordSessionInterface, channelID string) error {
		messageContent := "Sorry, no members found with this role"
		_, err := session.ChannelMessageSend(channelID, messageContent)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"channel_id": channelID,
				"error":      err,
			}).Errorf("Failed to 'no members found' message: %v", err)
			return err
		}
		logrus.WithField("channelID", channelID).Info("Successfully sent 'no members found' message.")
		return nil
	}

	handleDevModeFunc = func(session utils.DiscordSessionInterface, mentions []string, params CommandParams) error {
		logrus.WithFields(logrus.Fields{
			"channelID": params.ChannelID,
			"roleID":    params.RoleID,
			"mentions":  len(mentions),
			"batchSize": BatchSize,
		}).Info("Handling Dev Mode: Sending mentions in batches")

		if len(mentions) == 0 {
			logrus.Warn("No members found to mention in Dev Mode")
			return nil
		}

		var failedMentions []string
		var firstError error
		for i := 0; i < len(mentions); i += BatchSize {
			end := i + BatchSize
			if end > len(mentions) {
				end = len(mentions)
			}

			currentBatch := mentions[i:end]
			logrus.Debugf("processing batch: %d-%d", i, end-1)

			for _, mention := range currentBatch {
				msgContent := mention
				if params.Message != "" {
					msgContent = fmt.Sprintf("%s %s", params.Message, mention)
				}

				_, err := session.ChannelMessageSend(params.ChannelID, msgContent)

				if err != nil {
					logrus.WithFields(logrus.Fields{
						"channelID": params.ChannelID,
						"mention":   mention,
						"error":     err,
					}).Error("Failed to send individual mention in Dev Mode batch")

					failedMentions = append(failedMentions, mention)
					if firstError == nil {
						firstError = err
					}
					continue
				}

				logrus.Debugf("Successfully sent mention: %s", mention)
			}

			if end < len(mentions) {
				logrus.Infof("Rate limiting: sleeping for %v", BatchDelay)
				time.Sleep(BatchDelay)
			}
		}

		if len(failedMentions) > 0 {
			logrus.Warnf("Completed Dev Mode processing, but failed to send mentions to %d users: %v", len(failedMentions), failedMentions)
			summaryMsg := fmt.Sprintf("Finished mentioning, but failed for %d users: %s", len(failedMentions), strings.Join(failedMentions, ", "))
			_, _ = session.ChannelMessageSend(params.ChannelID, summaryMsg)
			return fmt.Errorf("failed to send %d out of %d mentions", len(failedMentions), len(mentions))
		} else {
			logrus.Infof("Dev Mode completed successfully")
		}

		return nil
	}

	handleDevTitleModeFunc = func(session utils.DiscordSessionInterface, mentions []string, params CommandParams) error {
		response := utils.FormatUserListResponse(mentions, params.RoleID)
		logrus.WithFields(logrus.Fields{
			"channelID":          params.ChannelID,
			"roleID":             params.RoleID,
			"response":           response,
			"GENERATED_RESPONSE": response,
		}).Info("Handling Dev Title Mode: Sending response")

		_, err := session.ChannelMessageSend(params.ChannelID, response)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"channelID": params.ChannelID,
				"roleID":    params.RoleID,
				"error":     err,
			}).Error("Failed to send dev_title response")
			return err
		}
		logrus.Infof("Successfully sent dev_title response")
		return nil
	}
	handleStandardModeFunc = func(session utils.DiscordSessionInterface, mentions []string, params CommandParams) error {

		if len(mentions) == 0 {
			logrus.Warnf("handleStandardModeFunc called with zero mentions for role %s. Sending 'no user' message.", params.RoleID)
			noUserMessage := "Sorry no user found under this role."
			_, err := session.ChannelMessageSend(params.ChannelID, noUserMessage)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"channelID": params.ChannelID,
					"roleID":    params.RoleID,
					"error":     err,
				}).Error("Failed to send 'no user found' message from handleStandardModeFunc")
				return err
			}
			return nil
		}
		response := utils.FormatMentionResponse(mentions, params.Message)
		logrus.WithFields(logrus.Fields{
			"channelID": params.ChannelID,
			"roleID":    params.RoleID,
			"response":  response,
		}).Info("Handling Standard Mode: Sending response")
		_, err := session.ChannelMessageSend(params.ChannelID, response)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"channelID": params.ChannelID,
				"roleID":    params.RoleID,
				"error":     err,
			}).Error("Failed to send mention response")
			return err
		}
		logrus.Infof("Successfully sent mention response")
		return nil
	}
)

func (s *CommandHandler) mentionEachHandler() error {
	logrus.Info("Processing mention-each command")

	params, err := extractCommandParamsFunc(s.discordMessage.MetaData)
	if err != nil {
		return fmt.Errorf("failed to extract command params: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"roleID":    params.RoleID,
		"channelID": params.ChannelID,
		"guildID":   params.GuildID,
		"message":   params.Message,
		"dev":       params.Dev,
		"devTitle":  params.DevTitle,
	}).Info("Extracted command parameters")

	discordSession, err := CreateSession()
	if err != nil {
		logrus.Error("Failed to create Discord session: ", err)
		return fmt.Errorf("failed to create Discord session: %w", err)
	}

	sessionWrapper := &DiscordSessionWrapper{discordSession}

	defer func() {
		if closeErr := discordSession.Close(); closeErr != nil {
			logrus.Errorf("Error closing session: %v", closeErr)
		}
	}()

	members, err := fetchMembersWithRoleFunc(sessionWrapper, params.GuildID, params.RoleID, params.ChannelID)
	if err != nil {
		return err
	}

	if len(members) == 0 {
		return sendNoMembersMessageFunc(sessionWrapper, params.ChannelID)
	}

	mentions := utils.FormatUserMentions(members)

	if params.DevTitle {
		return handleDevTitleModeFunc(sessionWrapper, mentions, params)
	} else if params.Dev {
		return handleDevModeFunc(sessionWrapper, mentions, params)
	} else {
		return handleStandardModeFunc(sessionWrapper, mentions, params)
	}
}
