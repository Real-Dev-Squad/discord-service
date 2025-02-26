package fixtures

import (
	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
)

var HelloCommand = dtos.DiscordMessage{
	AppPermissions: "0",
	ApplicationId:  "123456789012345678",
	Type:           discordgo.InteractionApplicationCommand,
	Channel: &discordgo.Channel{
		ID:   "987654321098765432",
		Name: "general",
	},
	ChannelId: "987654321098765432",
	Member: &discordgo.Member{
		User: &discordgo.User{
			ID:            "123456789012345678",
			Username:      "ExampleUser",
			Discriminator: "1234",
		},
	},
	Data: &dtos.Data{
		GuildId: "876543210987654321",
		ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
			Name: utils.CommandNames.Hello,
		},
	},
}
