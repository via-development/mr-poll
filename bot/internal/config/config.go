package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	BotToken     string
	BotId        string
	BotPublicKey string

	DSN         string
	AutoMigrate bool
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{}

	if config.BotToken = os.Getenv("BOT_TOKEN"); config.BotToken == "" {
		return nil, keyMissingError("BOT_TOKEN")
	}

	if config.BotId = os.Getenv("BOT_ID"); config.BotId == "" {
		return nil, keyMissingError("BOT_ID")
	}

	if config.BotPublicKey = os.Getenv("BOT_PUBLIC_KEY"); config.BotPublicKey == "" {
		return nil, keyMissingError("BOT_PUBLIC_KEY")
	}

	if config.DSN = os.Getenv("DSN"); config.DSN == "" {
		return nil, keyMissingError("DSN")
	}

	config.AutoMigrate = false
	
	return config, nil
}

func keyMissingError(key string) error {
	return errors.New(fmt.Sprintf("environment variable %s not set", key))
}
