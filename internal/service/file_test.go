package service

import (
	"bytes"
	"context"
	"mime/multipart"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/escoutdoor/social/internal/testutils"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type fileServiceSuite struct {
	suite.Suite
	container testcontainers.Container
	svc       File
}

func (st *fileServiceSuite) SetupSuite() {
	container, s3, err := testutils.NewMinIOContainer()
	st.Require().NoError(err, "failed to run minio container")
	st.Require().NotEmpty(container, "expected to get minio container")
	st.Require().NotEmpty(s3, "expected to get minio connection")

	st.container = container
	st.svc = NewFileService(s3)
}

func (st *fileServiceSuite) TearDownSuite() {
	err := st.container.Terminate(context.Background())
	st.Require().NoError(err, "failed to terminate postgres container")
}

func (st *fileServiceSuite) TestCreate() {
	var (
		body    bytes.Buffer
		content = []byte("wassup")
		ctx     = context.Background()
	)

	mw := multipart.NewWriter(&body)
	fw, err := mw.CreateFormFile("file", "file.jpg")
	st.NoError(err, "failed to create form file")

	n, err := fw.Write(content)
	st.NoError(err, "failed to write content into file")
	st.Equal(len(content), n)

	err = mw.Close()
	st.NoError(err, "failed to close writer")

	hdr := &multipart.FileHeader{
		Filename: gofakeit.BeerName(),
		Size:     int64(len(body.Bytes())),
	}

	url, err := st.svc.Create(ctx, &body, hdr)
	st.NoError(err, "failed to store photo into s3")
	st.NotEmpty(url, "expected to get photo url")
}

func TestFileService(t *testing.T) {
	suite.Run(t, new(fileServiceSuite))
}
