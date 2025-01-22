package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Port               string
	DISCORD_PUBLIC_KEY string
	GUILD_ID           string
	BOT_TOKEN          string
	QUEUE_NAME         string
	MAX_RETRIES        int
}

var AppConfig Config

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Error(err)
	} else {
		logrus.Info("Loaded .env file successfully")
	}

	AppConfig = Config{
		Port:               loadEnv("PORT"),
		DISCORD_PUBLIC_KEY: loadEnv("DISCORD_PUBLIC_KEY"),
		GUILD_ID:           loadEnv("GUILD_ID"),
		BOT_TOKEN:          loadEnv("BOT_TOKEN"),
		QUEUE_NAME:         "DISCORD_QUEUE",
		MAX_RETRIES:        5,
	}
}

func loadEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s not set", key))
	}
	return value
}
