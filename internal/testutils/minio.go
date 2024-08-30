package testutils

import (
	"context"
	"fmt"

	"github.com/escoutdoor/social/internal/s3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	minioUser = "test-user"
	minioPw   = "escoutdoor2024"
)

func NewMinIOContainer() (testcontainers.Container, s3.Repository, error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "minio/minio:latest",
		ExposedPorts: []string{"9000/tcp"},
		Env: map[string]string{
			"MINIO_ROOT_USER":     minioUser,
			"MINIO_ROOT_PASSWORD": minioPw,
		},
		Cmd:        []string{"server", "/data"},
		WaitingFor: wait.ForHTTP("/minio/health/live").WithPort("9000"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
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
		MinIOBucketName: "testbucket",
		MinIOEndpoint:   endpoint,
		MinIOHost:       endpoint,
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
