package constants

import "github.com/bwmarrin/discordgo"

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "hello",
		Description: "Greets back with hello!",
	},
}
