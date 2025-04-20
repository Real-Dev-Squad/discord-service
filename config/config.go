package config

import (
	"fmt"
	"os"
	"strconv"

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
	ENV                   Environment
	MAX_RETRIES           int
	RDS_BASE_API_URL      string
	VERIFICATION_SITE_URL string
	MENTION_EACH_ENABLED  bool
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
		MENTION_EACH_ENABLED: loadBoolEnv("MENTION_EACH_ENABLED", false),
	}

	// Loading Constants
	AppConfig.RDS_BASE_API_URL = EnvironmentURLs[AppConfig.ENV].RDS_BASE_API_URL
	AppConfig.VERIFICATION_SITE_URL = EnvironmentURLs[AppConfig.ENV].VERIFICATION_SITE_URL
	logrus.Infof("Feature Flag Status: MENTION_EACH_ENABLED=%t", AppConfig.MENTION_EACH_ENABLED)

}

func loadEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logrus.Panic(fmt.Sprintf("Environment variable %s not set", key))
	}
	return value
}

func loadBoolEnv(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		logrus.Infof("Environment variable %s not set, defaulting to false", key)
		return defaultValue
	}
	valueBool, err := strconv.ParseBool(valueStr)
	if err != nil {
		logrus.Warnf("Invalid boolean value for environment variable %s: '%s'. Defaulting to false.", key, valueStr)
		return defaultValue
	}
	logrus.Infof("Loaded %s=%t", key, valueBool)
	return valueBool
}