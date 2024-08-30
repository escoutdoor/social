package service

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/escoutdoor/social/internal/s3"
	"github.com/escoutdoor/social/internal/types"
)

type FileService struct {
	s3 s3.Repository
}

func NewFileService(s3 s3.Repository) *FileService {
	return &FileService{
		s3: s3,
	}
}

func (s *FileService) Create(ctx context.Context, src io.Reader, hdr *multipart.FileHeader) (string, error) {
	f := types.File{
		Name:    hdr.Filename,
		Payload: src,
		Size:    hdr.Size,
	}

	url, err := s.s3.Create(f)
	if err != nil {
		return "", err
	}
	return url, nil
}
