package constants

import "github.com/bwmarrin/discordgo"

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "hello",
		Description: "Greets back with hello!",
	},
	{
		Name:        "listening",
		Description: "mark user as listening",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "value",
				Description: "to enable or disable the listening mode",
				Type:        5,
				Required:    true,
			},
		},
	},
}
