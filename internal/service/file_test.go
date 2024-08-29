package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type fileServiceSuite struct {
	suite.Suite
	container      testcontainers.Container
	minioContainer testcontainers.Container
	svc            File
}

func (st *fileServiceSuite) SetupTest() {
	// container, db, err := testutils.NewPostgresContainer()
	// st.Require().NoError(err, "failed to run postgres container")
	// st.Require().NotEmpty(container, "expected to get postgres container")
	// st.Require().NotEmpty(db, "expected to get db connection")
	//
	// minioContainer, s3, err := testutils.NewMinIOContainer()
	// st.Require().NoError(err, "failed to run minio container")
	// st.Require().NotEmpty(minioContainer, "expected to get minio container")
	// st.Require().NotEmpty(s3, "expected to get minio connection")
	//
	// st.container = container
	// st.minioContainer = minioContainer
	// st.svc = NewFileService(s3)
}

func (st *fileServiceSuite) TestCreate() {
	// ctx := context.Background()

	// url, err := st.svc.Create(ctx, file, hdr)
	// st.Require().NoError(err, "failed to create file")
	// st.Require().NotEmpty(url, "expected to get file URL")
}

func TestFileService(t *testing.T) {
	suite.Run(t, new(fileServiceSuite))
}
