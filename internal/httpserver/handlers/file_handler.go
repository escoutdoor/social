package handlers

import (
	"log/slog"
	"net/http"

	"github.com/escoutdoor/social/internal/httpserver/responses"
	"github.com/escoutdoor/social/internal/service"
	"github.com/go-chi/chi/v5"
)

type FileHandler struct {
	svc service.File
}

func NewFileHandler(svc service.File) FileHandler {
	return FileHandler{
		svc: svc,
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

	ctx := r.Context()
	url, err := h.svc.Create(ctx, src, hdr)
	if err != nil {
		slog.Error("FileHandler.Create - FileService.Create", "error", err)
		responses.InternalServerResponse(w, ErrFileSaveFailed)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "file successfully uploaded", "url": url})
}
