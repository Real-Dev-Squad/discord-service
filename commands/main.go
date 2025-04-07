package constants

import (
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        utils.CommandNames.Hello,
		Description: "Greets back with hello!",
	},
	{
		Name:        utils.CommandNames.Listening,
		Description: "Mark user as listening",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "value",
				Description: "to enable or disable the listening mode",
				Type:        5,
				Required:    true,
			},
		},
	},
	{
		Name:        utils.CommandNames.Verify,
		Description: "Generate a link with user specific token to link with RDS backend",
	},
	{
		Name:        utils.CommandNames.MentionEach,
		Description: "Mention all users with a specific role",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "role",
				Description: "The role whose members will be mentioned",
				Type:        discordgo.ApplicationCommandOptionRole,
				Required:    true,
			},
			{
				Name:        "message",
				Description: "Custom message to accompany the mentions",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
			{
				Name:        "dev",
				Description: "Send individual mentions in separate message",
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Required:    false,
			},
			{
				Name:        "dev_title",
				Description: "Show a list of users with the role without pinging them",
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Required:    false,
			},
		},
	},
}
