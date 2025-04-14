package utils

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"strings"
)

type DiscordSessionInterface interface {
	GuildMembers(guildID, after string, limit int) ([]*discordgo.Member, error)
	ChannelMessageSend(channelID, content string) (*discordgo.Message, error)
}

func GetUsersWithRole(session DiscordSessionInterface, guildID string, roleID string) ([]*discordgo.Member, error) {
	var membersWithRole []*discordgo.Member
	lastMemberID := ""
	limit := 1000
	for {
		membersChunk, err := session.GuildMembers(guildID, lastMemberID, limit)
		if err != nil {
			logrus.Errorf("failed to fetch guild members chunk for guild %s: %v", guildID, err)
			return nil, fmt.Errorf("failed to fetch guild members chunk: %w", err)
		}
		if len(membersChunk) == 0 {
			break
		}
		foundInChunk := 0
		for _, member := range membersChunk {
			if member == nil || member.User == nil {
				logrus.Warnf("Guild %s: Member object or User data is nil, skipping member: %+v", guildID, member)
				continue
			}
			if HasRole(member, roleID) {
				membersWithRole = append(membersWithRole, member)
				foundInChunk++
			}
			lastMemberID = member.User.ID
		}
		if len(membersChunk) < limit {
			break
		}
	}

	logrus.Infof("Finished fetching. Found %d total members with role %s in guild %s", len(membersWithRole), roleID, guildID)
	return membersWithRole, nil
}

func HasRole(member *discordgo.Member, roleID string) bool {
	if member == nil || member.Roles == nil {
		return false
	}

	for _, r := range member.Roles {
		if r == roleID {
			return true
		}
	}

	return false
}

func FormatUserMentions(members []*discordgo.Member) []string {
	if members == nil || len(members) == 0 {
		return []string{}
	}
	mentions := make([]string, len(members))
	for i, member := range members {
		if member != nil && member.User != nil {
			mentions[i] = fmt.Sprintf("<@%s>", member.User.ID)
		} else {
			logrus.Warnf("Attempted to format mention for nil member or user at index %d", i)
			mentions[i] = "[invalid user data]"
		}
	}
	return mentions
}

func FormatRoleMention(roleID string) string {
	return fmt.Sprintf("<@&%s>", roleID)
}

func FormatMentionResponse(mentions []string, message string) string {
	if len(mentions) == 0 {
		return "Sorry no user found under this role."
	}

	mentionStrings := strings.Join(mentions, " ")
	if message == "" {
		return mentionStrings
	}

	return fmt.Sprintf("%s %s", message, mentionStrings)

}

func FormatDevTitleResponse(mentions []string, roleID string) string {
	count := len(mentions)
	roleMention := FormatRoleMention(roleID)

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
