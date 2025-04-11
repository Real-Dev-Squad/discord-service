package utils

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"slices"
	"strings"
)

type DiscordSessionInterface interface {
	GuildMembers(guildID, after string, limit int) ([]*discordgo.Member, error)
	ChannelMessageSend(channelID, content string) (*discordgo.Message, error)
}

func GetUsersWithRole(session DiscordSessionInterface, guildID string, roleID string) ([]*discordgo.Member, error) {
	members, err := session.GuildMembers(guildID, "", 1000)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch guild members: %w", err)
	}

	var membersWithRole []*discordgo.Member
	for _, member := range members {
		if slices.Contains(member.Roles, roleID) {
			membersWithRole = append(membersWithRole, member)
		}
	}

	return membersWithRole, nil
}

func FormatUserMentions(members []*discordgo.Member) []string {
	if members == nil {
		return []string{}
	}
	mentions := make([]string, len(members))
	for i, member := range members {
		mentions[i] = fmt.Sprintf("<@%s>", member.User.ID)
	}
	return mentions
}

func FormatRoleMention(roleID string) string {
	return fmt.Sprintf("<@&%s>", roleID)
}

func JoinMentions(mentions []string, separator string) string {
	if len(mentions) == 0 {
		return ""
	}
	return strings.Join(mentions, separator)
}

func FormatMentionResponse(mentions []string, message string) string {
	if len(mentions) == 0 {
		return "Sorry no user found under this role."
	}

	if message == "" {
		return JoinMentions(mentions, " ")
	}

	return fmt.Sprintf("%s %s", message, JoinMentions(mentions, " "))
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
		userList := JoinMentions(mentions, ", ")
		return fmt.Sprintf("Found %d users with the %s role: %s", count, roleMention, userList)
	}
}
