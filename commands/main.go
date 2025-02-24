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
}
