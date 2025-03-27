package utils

import (
	"fmt"
	"slices"
	"github.com/bwmarrin/discordgo"
)

type DiscordSessionInterface interface {
	GuildMembers(guildID, after string, limit int) ([]*discordgo.Member, error)
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
	if mentions == nil || len(mentions) == 0 {
		return ""
	}

	results := mentions[0]
	for i := 1; i < len(mentions); i++ {
		results += separator + mentions[i]
	}

	return results
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
	if mentions == nil || len(mentions) == 0 {
		return fmt.Sprintf("Sorry, no user found with <@&%s> role.", roleID)
	} else if len(mentions) == 1 {
		return fmt.Sprintf("The user with <@&%s> role is %s.", roleID, mentions[0])
	} else {
		return fmt.Sprintf("The users with <@&%s> role are %s.", roleID, JoinMentions(mentions, ", "))
	}
}