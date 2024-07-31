package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port        int    `envconfig:"PORT" default:"8080"`
	PostgresURL string `envconfig:"POSTGRES_URL" required:"true"`
	JWTKey      string `envconfig:"JWT_SIGN_KEY" required:"true"`

	MinIOHost       string `envconfig:"MINIO_HOST"`
	MinIOEndpoint   string `envconfig:"MINIO_SERVER_URL"`
	MinIOUser       string `envconfig:"MINIO_ROOT_USER"`
	MinIOPw         string `envconfig:"MINIO_ROOT_PASSWORD"`
	MinIOBucketName string `envconfig:"MINIO_BUCKET_NAME"`
	MinIOUseSSL     bool   `envconfig:"MINIO_USE_SSL" default:"false"`
	MinIORegion     string `envconfig:"MINIO_REGION" default:"auto"`
}

func New() (Config, error) {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}
	fmt.Printf("config: %+v\n", cfg)
	return cfg, nil
}
