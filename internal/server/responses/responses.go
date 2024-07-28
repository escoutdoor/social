package responses

import (
	"encoding/json"
	"net/http"
)

type ErrResponse struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrResponse{Error: message})
}

func InternalServerResponse(w http.ResponseWriter, err error) {
	ErrorResponse(w, http.StatusInternalServerError, err.Error())
}

func NotFoundResponse(w http.ResponseWriter, err error) {
	ErrorResponse(w, http.StatusNotFound, err.Error())
}

func BadRequestResponse(w http.ResponseWriter, err error) {
	ErrorResponse(w, http.StatusBadRequest, err.Error())
}

func UnauthorizedResponse(w http.ResponseWriter, err error) {
	ErrorResponse(w, http.StatusUnauthorized, err.Error())
}

func ForbiddenResponse(w http.ResponseWriter, err error) {
	ErrorResponse(w, http.StatusForbidden, err.Error())
}
