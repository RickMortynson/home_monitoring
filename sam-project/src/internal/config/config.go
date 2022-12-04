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
	SupabaseBucket string `env:"SUPABASE_BUCKET"`
	SupabaseToken  string `env:"SUPABASE_SERVICE_TOKEN"`
	MetadataFile   string `env:"METADATA_FILENAME"`
}

func LoadConfig(configs ...any) error {
	var envFile string
	if os.Getenv("GO_ENV") == "production" {
		envFile = "secrets/.prod.env"
	} else {
		envFile = "secrets/.dev.env"
	}

	for _, cfg := range configs {
		if err := cleanenv.ReadConfig(envFile, cfg); err != nil {
			return err
		}
	}

	return nil
}
