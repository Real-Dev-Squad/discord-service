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

func (sw *SessionWrapper) GetUerId() string {
	return sw.Session.State.User.ID
}

type SessionInterface interface {
	Open() error
	Close() error
	ApplicationCommandCreate(applicationID, guildID string, command *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error)
	GetUerId() string
}
