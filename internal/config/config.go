package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env         string `envconfig:"ENV" default:"local"`
	Port        int    `envconfig:"PORT" default:"8080"`
	PostgresURL string `envconfig:"POSTGRES_URL" required:"true"`
	RedisURL    string `envconfig:"REDIS_URL" required:"true"`
	SignKey     string `envconfig:"JWT_SIGN_KEY" required:"true"`

	MinIOHost       string `envconfig:"MINIO_HOST" required:"true"`
	MinIOEndpoint   string `envconfig:"MINIO_SERVER_URL" required:"true"`
	MinIOUser       string `envconfig:"MINIO_ROOT_USER" required:"true"`
	MinIOPw         string `envconfig:"MINIO_ROOT_PASSWORD" required:"true"`
	MinIOBucketName string `envconfig:"MINIO_BUCKET_NAME" required:"true"`
	MinIOUseSSL     bool   `envconfig:"MINIO_USE_SSL" default:"false"`
	MinIORegion     string `envconfig:"MINIO_REGION" default:"auto"`
}

func New() (*Config, error) {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return &cfg, err
	}
	return &cfg, nil
}
