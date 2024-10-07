package dtos

import "github.com/bwmarrin/discordgo"

type DiscordMessage struct {
	Type discordgo.InteractionType `json:"type"`
	User *discordgo.User           `json:"user"`
}
