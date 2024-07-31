package s3

import (
	"context"
	"fmt"

	"github.com/escoutdoor/social/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	mc  *minio.Client
	cfg config.Config
}

func New(cfg config.Config) (*MinIOClient, error) {
	ctx := context.Background()
	client, err := minio.New(cfg.MinIOHost, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOUser, cfg.MinIOPw, ""),
		Secure: cfg.MinIOUseSSL,
		Region: cfg.MinIORegion,
	})
	if err != nil {
		return nil, err
	}

	ie, err := client.BucketExists(ctx, cfg.MinIOBucketName)
	if err != nil {
		return nil, err
	}
	if !ie {
		if err := client.MakeBucket(ctx, cfg.MinIOBucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
		if err := client.SetBucketPolicy(
			ctx,
			cfg.MinIOBucketName,
			`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Principal": {
							"AWS": [
								"*"
							]
						},
						"Action": [
							"s3:GetObject"
						],
						"Resource": [
							"arn:aws:s3:::`+cfg.MinIOBucketName+`/*"
						]
					}
				]
			}`); err != nil {
			return nil, fmt.Errorf("error client.SetBucketPolicy: %w", err)
		}
	}
	return &MinIOClient{mc: client, cfg: cfg}, nil
}
