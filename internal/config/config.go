package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	JWTKey      string `envconfig:"JWT_SIGN_KEY" required:"true"`
	PostgresURL string `envconfig:"POSTGRES_URL" required:"true"`
	Port        int    `envconfig:"PORT" default:"8080"`
}

func New() (Config, error) {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	fmt.Printf("config data: %+v\n", cfg)
	return cfg, nil
}
