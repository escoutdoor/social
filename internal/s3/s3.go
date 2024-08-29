package s3

import (
	"context"
	"fmt"

	"github.com/escoutdoor/social/internal/types"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	mc *minio.Client
	Opts
}

type Opts struct {
	MinIOBucketName string
	MinIOEndpoint   string
	MinIOHost       string
	MinIOUser       string
	MinIOPw         string
	MinIOUseSSL     bool
	MinIORegion     string
}

type Repository interface {
	Create(file types.File) (string, error)
	Delete(id string) error
	GetByID(id string) (string, error)
}

func New(opts Opts) (*MinIOClient, error) {
	ctx := context.Background()
	client, err := minio.New(opts.MinIOHost, &minio.Options{
		Creds:  credentials.NewStaticV4(opts.MinIOUser, opts.MinIOPw, ""),
		Secure: opts.MinIOUseSSL,
		Region: opts.MinIORegion,
	})
	if err != nil {
		return nil, err
	}

	ie, err := client.BucketExists(ctx, opts.MinIOBucketName)
	if err != nil {
		return nil, err
	}
	if !ie {
		if err := client.MakeBucket(ctx, opts.MinIOBucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
		if err := client.SetBucketPolicy(
			ctx,
			opts.MinIOBucketName,
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
							"arn:aws:s3:::`+opts.MinIOBucketName+`/*"
						]
					}
				]
			}`); err != nil {
			return nil, fmt.Errorf("error client.SetBucketPolicy: %w", err)
		}
	}
	return &MinIOClient{mc: client}, nil
}
