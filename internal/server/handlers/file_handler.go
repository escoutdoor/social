package handlers

import (
	"log/slog"
	"net/http"

	"github.com/escoutdoor/social/internal/s3"
	"github.com/escoutdoor/social/internal/server/responses"
	"github.com/escoutdoor/social/internal/types"
	"github.com/go-chi/chi"
)

type FileHandler struct {
	minio *s3.MinIOClient
}

func NewFileHandler(m *s3.MinIOClient) FileHandler {
	return FileHandler{
		minio: m,
	}
}

func (h *FileHandler) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", h.create)
	return r
}

func (h *FileHandler) create(w http.ResponseWriter, r *http.Request) {
	src, hdr, err := r.FormFile("file")
	if err != nil {
		responses.BadRequestResponse(w, ErrFileNotReceived)
		return
	}
	defer src.Close()
	f := types.File{
		Name:    hdr.Filename,
		Payload: src,
		Size:    hdr.Size,
	}

	url, err := h.minio.Create(f)
	if err != nil {
		slog.Error("FileHandler.Create - MinIOClient.Create", "error", err)
		responses.InternalServerResponse(w, ErrFileSaveFailed)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "file successfully uploaded", "url": url})
}
