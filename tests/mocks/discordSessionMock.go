package mocks

import (
	// Import the actual interface definition

	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
)

type DiscordSession struct {
	mock.Mock
}

// Implement models.SessionInterface
func (m *DiscordSession) Open() error {
	args := m.Called()
	return args.Error(0)
}

func (m *DiscordSession) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *DiscordSession) ApplicationCommandCreate(appID, guildID string, cmd *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error) {
	args := m.Called(appID, guildID, cmd)
	var retCmd *discordgo.ApplicationCommand
	if args.Get(0) != nil {
		var ok bool
		retCmd, ok = args.Get(0).(*discordgo.ApplicationCommand)
		if !ok {
			panic(fmt.Sprintf("mock return value 0 for ApplicationCommandCreate is not *discordgo.ApplicationCommand: %T", args.Get(0)))
		}
	}
	return retCmd, args.Error(1)
}

func (m *DiscordSession) GetUserId() string {
	args := m.Called()
	return args.String(0)
}

func (m *DiscordSession) GuildMembers(guildID, after string, limit int) ([]*discordgo.Member, error) {
	args := m.Called(guildID, after, limit)
	var retMembers []*discordgo.Member
	if args.Get(0) != nil {
		var ok bool
		retMembers, ok = args.Get(0).([]*discordgo.Member)
		if !ok {
			panic(fmt.Sprintf("mock return value 0 for GuildMembers is not []*discordgo.Member: %T", args.Get(0)))
		}
	}
	return retMembers, args.Error(1)
}

func (m *DiscordSession) ChannelMessageSend(channelID, content string) (*discordgo.Message, error) {
	args := m.Called(channelID, content)
	var retMsg *discordgo.Message
	if args.Get(0) != nil {
		var ok bool
		retMsg, ok = args.Get(0).(*discordgo.Message)
		if !ok {
			panic(fmt.Sprintf("mock return value 0 for ChannelMessageSend is not *discordgo.Message: %T", args.Get(0)))
		}
	}
	return retMsg, args.Error(1)
}
