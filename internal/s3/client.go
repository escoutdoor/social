package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

const expires = time.Second * 24 * 60 * 60

func (m *MinIOClient) generateUrl(id string) string {
	url := fmt.Sprintf("%s/%s/%s", m.MinIOEndpoint, m.MinIOBucketName, id)
	return url
}

func (m *MinIOClient) Create(file types.File) (string, error) {
	id := uuid.New().String()
	if _, err := m.mc.PutObject(
		context.Background(),
		m.MinIOBucketName,
		id,
		file.Payload,
		file.Size,
		minio.PutObjectOptions{ContentType: "image/png"},
	); err != nil {
		return "", err
	}

	// pr, err := m.mc.PresignedGetObject(context.Background(), m.cfg.MinIOBucketName, id, expires, nil)
	// if err != nil {
	// 	return "", err
	// }
	url := m.generateUrl(id)
	return url, nil
}

func (m *MinIOClient) GetByID(id string) (string, error) {
	// url, err := m.mc.PresignedGetObject(context.Background(), m.cfg.MinIOBucketName, id, expires, nil)
	// if err != nil {
	// 	return "", err
	// }
	// return url.String(), nil

	url := m.generateUrl(id)
	return url, nil
}

func (m *MinIOClient) Delete(id string) error {
	err := m.mc.RemoveObject(context.Background(), m.MinIOBucketName, id, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
