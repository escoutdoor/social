package testutils

import (
	"context"
	"fmt"

	"github.com/escoutdoor/social/internal/s3"
	"github.com/testcontainers/testcontainers-go"
	miniomodule "github.com/testcontainers/testcontainers-go/modules/minio"
)

var (
	minioUser = "test-user"
	minioPw   = "test-pw"
)

func NewMinIOContainer() (testcontainers.Container, s3.Repository, error) {
	ctx := context.Background()

	container, err := miniomodule.Run(ctx, "minio/minio:latest")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start MinIO container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get MinIO container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "9000")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get MinIO container port: %w", err)
	}

	endpoint := fmt.Sprintf("%s:%s", host, port.Port())
	s3, err := s3.New(s3.Opts{
		MinIOBucketName: "test-bucket",
		MinIOEndpoint:   endpoint,
		MinIOHost:       host,
		MinIOUser:       minioUser,
		MinIOPw:         minioPw,
		MinIOUseSSL:     false,
		MinIORegion:     "auto",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to s3: %w", err)
	}

	return container, s3, nil
}
