package constants

import "github.com/bwmarrin/discordgo"

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "hello",
		Description: "Greets back with hello!",
	},
	{
		Name:        "listening",
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
		Name:        "verify",
		Description: "Generate a link with user specific token to link with RDS backend",
	},
}
