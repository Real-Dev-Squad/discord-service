package utils

import "github.com/Real-Dev-Squad/discord-service/dtos"

const (
	NICKNAME_SUFFIX                = "-Can't Talk"
	NICKNAME_PREFIX                = "🎧 "
	DISCORD_GUILD_MEMBER_API_LIMIT = 1000
)

var CommandNames = dtos.CommandNameTypes{
	Hello:       "hello",
	Listening:   "listening",
	Verify:      "verify",
	MentionEach: "mention-each",
}
