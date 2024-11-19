package dtos

import "github.com/bwmarrin/discordgo"

type data struct {
	discordgo.ApplicationCommandInteractionData
	GuildId string `json:"guild_id"`
}

type DiscordMessage struct {
	AppPermissions string                    `json:"app_permissions"`
	ApplicationId  string                    `json:"application_id"`
	Type           discordgo.InteractionType `json:"type"`
	Channel        *discordgo.Channel        `json:"channel"`
	ChannelId      string                    `json:"channel_id"`
	Member         *discordgo.Member         `json:"member"`
	Data           *data                     `json:"data"`
}
