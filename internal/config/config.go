package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresURL string `envconfig:"POSTGRES_URL"`
}

func New() (Config, error) {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	fmt.Printf("config data: %+v\n", cfg)
	return cfg, nil
}
