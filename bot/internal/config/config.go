package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken     string
	BotId        string
	BotPublicKey string

	DSN         string
	AutoMigrate bool

	SentryDSN string

	ShardIds   []int
	ShardCount int

	EmbedColor int

	WebsiteURL string
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

	config.SentryDSN = os.Getenv("SENTRY_DSN")

	sc := os.Getenv("SHARD_COUNT")
	if sc != "" {
		config.ShardCount, err = strconv.Atoi(sc)
		if err != nil {
			return nil, err
		}
	} else {
		config.ShardCount = 1
	}

	si := os.Getenv("SHARD_IDS")
	if si != "" {
		for _, i := range strings.Split(si, ",") {
			id, err := strconv.Atoi(i)
			if err != nil {
				return nil, err
			}
			config.ShardIds = append(config.ShardIds, id)
		}
	} else {
		config.ShardIds = []int{0}
	}

	c := os.Getenv("EMBED_COLOR")
	if sc != "" {
		e, err := strconv.ParseInt(c, 16, 16)
		if err != nil {
			return nil, err
		}
		config.EmbedColor = int(e)
	} else {
		config.EmbedColor = 0x40FFAC
	}

	u := os.Getenv("WEBSITE_URL")
	if u != "" {
		config.WebsiteURL = u
	} else {
		return nil, keyMissingError("WEBSITE_URL")
	}

	am := os.Getenv("AUTO_MIGRATE")
	if am == "t" || am == "true" || am == "1" {
		config.AutoMigrate = true
	}

	return config, nil
}

func keyMissingError(key string) error {
	return errors.New(fmt.Sprintf("environment variable %s not set", key))
}
