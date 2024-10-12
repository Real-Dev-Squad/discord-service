package helpers

import "os"

func init() {
	os.Setenv("MODE", "test")
	os.Setenv("PORT", "8080")
	os.Setenv("DISCORD_PUBLIC_KEY", "<DISCORD_PUBLIC_KEY>")
}
