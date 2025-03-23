package tests

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

type MockSession struct {
	OpenError                        error
	CommandError                     error
	ApplicationCommandCalled         bool
	CloseCalled                      bool
	GetUserIdCalled                  bool
	ForceErrorForGuildMemberNickname bool
}

func (m *MockSession) Open() error {
	return m.OpenError
}

func (m *MockSession) Close() error {
	m.CloseCalled = true
	return nil
}

func (m *MockSession) ApplicationCommandCreate(applicationID, guildID string, command *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error) {
	if m.CommandError == nil {
		m.ApplicationCommandCalled = true
	}
	return nil, m.CommandError
}

func (m *MockSession) GetUserId() string {
	m.GetUserIdCalled = true
	return ""
}

func (m *MockSession) GuildMemberNickname(userId string, newNickName string) error {
	if m.ForceErrorForGuildMemberNickname {
		return assert.AnError
	}
	return nil
}
