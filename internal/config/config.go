package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	DB struct {
		User string `env:"DB_USER" env-default:"default"`
		Pass string `env:"DB_PASS" env-default:""`
		Host string `env:"DB_HOST" env-default:"localhost"`
		Port int    `env:"DB_PORT" env-default:"19000"`
		Name string `env:"DB_NAME" env-default:"default"`
	}

	Nats struct {
		URL    string `env:"NATS_URL" env-required:"true"`
		Queues struct {
			Validation string `env:"NATS_QUEUE_VALIDATION" env-required:"true"`
			Errors     string `env:"NATS_QUEUE_ERRORS" env-required:"true"`
			Status     string `env:"NATS_QUEUE_STATUS" env-required:"true"`
		}
	}
}

func New() *Config {
	config := &Config{}

	if err := cleanenv.ReadEnv(config); err != nil {
		header := "PRICES SERVICE ENVs"
		f := cleanenv.FUsage(os.Stdout, config, &header)
		f()
		panic(err)
	}

	return config
}
