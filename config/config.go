package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	Port                  string
	DISCORD_PUBLIC_KEY    string
	GUILD_ID              string
	BOT_TOKEN             string
	QUEUE_URL             string
	QUEUE_NAME            string
	ENV                   Environment
	MAX_RETRIES           int
	RDS_BASE_API_URL      string
	VERIFICATION_SITE_URL string
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
		ENV:                Environment(loadEnv("ENV")).Validate(),
		MAX_RETRIES:        5,
	}

	// Loading Constants
	AppConfig.RDS_BASE_API_URL = EnvironmentURLs[AppConfig.ENV].RDS_BASE_API_URL
	AppConfig.VERIFICATION_SITE_URL = EnvironmentURLs[AppConfig.ENV].VERIFICATION_SITE_URL

}

func loadEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logrus.Panic(fmt.Sprintf("Environment variable %s not set", key))
	}
	return value
}
