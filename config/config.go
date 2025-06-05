package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Port                  string
	DISCORD_PUBLIC_KEY    string
	GUILD_ID              string
	BOT_TOKEN             string
	QUEUE_URL             string
	QUEUE_NAME            string
	MAX_RETRIES           int
	RDS_BASE_API_URL      string
	MAIN_SITE_URL         string
	BOT_PRIVATE_KEY       string
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
		QUEUE_URL:          loadEnv("QUEUE_URL"),
		DISCORD_PUBLIC_KEY: loadEnv("DISCORD_PUBLIC_KEY"),
		GUILD_ID:           loadEnv("GUILD_ID"),
		BOT_TOKEN:          loadEnv("BOT_TOKEN"),
		QUEUE_NAME:         loadEnv("QUEUE_NAME"),
		MAX_RETRIES:        5,
		RDS_BASE_API_URL:   loadEnv("RDS_BASE_API_URL"),
		MAIN_SITE_URL:      loadEnv("MAIN_SITE_URL"),
		BOT_PRIVATE_KEY:    loadEnv("BOT_PRIVATE_KEY"),
	}
}

func loadEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logrus.Panic(fmt.Sprintf("Environment variable %s not set", key))
	}
	return value
}
