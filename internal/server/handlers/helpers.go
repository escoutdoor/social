package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func getIDParam(r *http.Request) (int, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return int(id), nil
}
