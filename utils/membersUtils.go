package utils

import (
	"fmt"
	"strings"

	"github.com/Real-Dev-Squad/discord-service/models"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func GetUsersWithRole(session models.SessionInterface, guildID string, roleID string) ([]discordgo.Member, error) {
	var membersWithRole []discordgo.Member
	lastMemberID := ""

	logrus.Debugf("Fetching members with role %s in guild %s", roleID, guildID)

	for {
		logrus.Debugf("Fetching members chunk after user ID: '%s', limit: '%d'", lastMemberID, DISCORD_GUILD_MEMBER_API_LIMIT)

		membersChunk, err := session.GuildMembers(guildID, lastMemberID, DISCORD_GUILD_MEMBER_API_LIMIT)
		if err != nil {
			logrus.Errorf("failed to fetch guild members chunk for guild %s: %v", guildID, err)
			return nil, fmt.Errorf("failed to fetch guild members chunk: %w", err)
		}
		logrus.Debugf("Fetched %d members in this chunk for guild %s", len(membersChunk), guildID)

		if len(membersChunk) == 0 {
			break
		}

		foundInChunk := 0
		lastIdInChunk := ""
		for _, member := range membersChunk {
			if member.User == nil {
				logrus.Warnf("Guild %s: Member object or User data is nil, skipping member: %+v", guildID, member)
				continue
			}
			for _, r := range member.Roles {
				if r == roleID {
					membersWithRole = append(membersWithRole, *member)
					foundInChunk++
					break
				}
			}

			lastIdInChunk = member.User.ID
		}
		lastMemberID = lastIdInChunk
		logrus.Debugf("Found %d members with role %s in this chunk (Guild %s)", foundInChunk, roleID, guildID)
	}

	logrus.Infof("Finished fetching. Found %d total members with role %s in guild %s", len(membersWithRole), roleID, guildID)
	return membersWithRole, nil
}

func FormatUserMentions(members []discordgo.Member) []string {
	mentions := make([]string, 0, len(members))
	for i, member := range members {
		if member.User != nil {
			mentions = append(mentions, fmt.Sprintf("<@%s>", member.User.ID))
		} else {
			logrus.Warnf("Skipping formatting mention for nil member or user at input index %d", i)
		}
	}
	return mentions
}

func FormatMentionResponse(mentions []string, message string) string {
	mentionStrings := strings.Join(mentions, " ")
	if message == "" {
		return mentionStrings
	}

	return fmt.Sprintf("%s %s", message, mentionStrings)

}

func FormatUserListResponse(mentions []string, roleID string) string {
	count := len(mentions)
	roleMention := fmt.Sprintf("<@&%s>", roleID)

	switch count {
	case 0:
		return fmt.Sprintf("Found 0 users with the %s role", roleMention)
	case 1:
		return fmt.Sprintf("Found 1 user with the %s role: %s", roleMention, mentions[0])
	default:
		userList := strings.Join(mentions, ", ")
		return fmt.Sprintf("Found %d users with the %s role: %s", count, roleMention, userList)
	}
}
