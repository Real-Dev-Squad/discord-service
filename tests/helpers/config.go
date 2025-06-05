package helpers

import (
	"os"
)

func init() {
	os.Setenv("PORT", "8080")
	os.Setenv("DISCORD_PUBLIC_KEY", "8933e3749b4feb4d76169b26ed372af3c378f4353c2024fee0601f2a2e7918e1")
	os.Setenv("GUILD_ID", "8933e3749b4feb4d76169b26ed372af3c378f4353c2024fee0601f2a2e7918e1")
	os.Setenv("BOT_TOKEN", "8933e3749b4feb4d76169b26ed372af3c378f4353c2024fee0601f2a2e7918e1")
	os.Setenv("QUEUE_NAME", "DISCORD_QUEUE")
	os.Setenv("QUEUE_URL", "local:5672")
	os.Setenv("ENV", "development")
	os.Setenv("RDS_BASE_API_URL", "http://localhost:3000")
	os.Setenv("MAIN_SITE_URL", "http://localhost:4200")
	os.Setenv("BOT_PRIVATE_KEY", "<discord-bot-private-key>")
}
