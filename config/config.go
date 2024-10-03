package config

import (
	"fmt"
	"os"

	utility "github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/joho/godotenv"
)

var logger = &utility.Logger{}

type PredefinedConfig struct {
	Port string
	Mode string
}

var Config PredefinedConfig

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	} else {
		logger.Info("Loaded .env file successfully")
	}

	Config.Mode = loadEnv("MODE")
	Config.Port = loadEnv("PORT")
}

func loadEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s not set", key))
	}
	return value
}
