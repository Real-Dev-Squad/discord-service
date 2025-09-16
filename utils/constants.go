package utils

import "github.com/Real-Dev-Squad/discord-service/dtos"

const NICKNAME_SUFFIX = "-Can't Talk"
const NICKNAME_PREFIX = "🎧 "

var CommandNames = dtos.CommandNameTypes{
	Hello:     "hello",
	Listening: "listening",
	Verify:    "verify",
}
