package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HeartbeatConfig struct {
	HeartbeatUrl    string `env:"HEARTBEAT_URL"`
	HeartbeatName   string `env:"OK_HEARTBEAT_NAME"`
	CheckTimeoutSec int    `env:"CHECK_TIMEOUT_SEC"`
}

type TelegramConfig struct {
	TgChatId   int    `env:"TG_CHAT_ID"`
	TgBotToken string `env:"TG_BOT_TOKEN"`
}

type RepoConfig struct {
	SupabaseUrl    string `env:"SUPABASE_URL"`
	SupabaseToken  string `env:"SUPABASE_SERVICE_TOKEN"`
}

func LoadConfig(configs ...any) error {
	var envFile string
	// doesn't work actually, lol
	// TODO: find out & fuck around with SAM's in-container env variables
	if os.Getenv("GO_ENV") == "production" {
		envFile = "secrets/.prod.env"
	} else {
		envFile = "secrets/.dev.env"
	}

	// TODO: get rid of this 🤦
	envFile = "secrets/.prod.env"

	for _, cfg := range configs {
		if err := cleanenv.ReadConfig(envFile, cfg); err != nil {
			return err
		}
	}

	return nil
}
