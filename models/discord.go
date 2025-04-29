package models

import "github.com/bwmarrin/discordgo"

type SessionWrapper struct {
	Session *discordgo.Session
}

func (s *SessionWrapper) Open() error {
	return s.Session.Open()
}

func (s *SessionWrapper) Close() error {
	return s.Session.Close()
}

func (s *SessionWrapper) ApplicationCommandCreate(applicationID, guildID string, command *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error) {
	return s.Session.ApplicationCommandCreate(applicationID, guildID, command)
}

func (sw *SessionWrapper) GetUserId() string {
	return sw.Session.State.User.ID
}

func (sw *SessionWrapper) GuildMembers(guildID, after string, limit int) ([]*discordgo.Member, error) {
	return sw.Session.GuildMembers(guildID, after, limit)
}

func (sw *SessionWrapper) ChannelMessageSend(channelID, content string) (*discordgo.Message, error) {
	return sw.Session.ChannelMessageSend(channelID, content)
}

type SessionInterface interface {
	Open() error
	Close() error
	ApplicationCommandCreate(applicationID, guildID string, command *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error)
	GetUserId() string
	GuildMembers(guildID, after string, limit int) ([]*discordgo.Member, error)
	ChannelMessageSend(channelID, content string) (*discordgo.Message, error)
}
