package utils

import "github.com/Real-Dev-Squad/discord-service/dtos"

const (
	NICKNAME_SUFFIX             = "-Can't Talk"
	NICKNAME_PREFIX             = "🎧 "
	DiscordGuildMembersAPILimit = 1000
)

var CommandNames = dtos.CommandNameTypes{
	Hello:       "hello",
	Listening:   "listening",
	Verify:      "verify",
	MentionEach: "mention-each",
}
