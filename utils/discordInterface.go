package utils

import "github.com/bwmarrin/discordgo"

type DiscordSessionInterface interface {
	GuildMembers(guildID, after string, limit int) ([]*discordgo.Member, error)
	ChannelMessageSend(channelID, content string) (*discordgo.Message, error)
}
