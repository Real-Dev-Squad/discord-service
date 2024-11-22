package utils

import "fmt"

type responseGenerator struct{}

func (r *responseGenerator) HelloResponse(userId string) string {
	return fmt.Sprintf("Hey there <@%s>! Congratulations, you just executed your first slash command", userId)
}

var ResponseGenerator = &responseGenerator{}
