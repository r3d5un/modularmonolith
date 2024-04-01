package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/r3d5un/modularmonolith/internal/logging"
)

func ReadStringParameter(key string, r *http.Request) (*string, error) {
	s := r.PathValue(key)
	if s == "" {
		return nil, fmt.Errorf("empty string parameter")
	}

	return &s, nil
}

type ErrorMessage struct {
	Message any `json:"message"`
}

func logError(r *http.Request, err error) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.ErrorContext(
		ctx,
		"an error occurred",
		"request_method", r.Method,
		"request_url", r.URL.String(),
		"error", err,
	)
}

func errorResponse(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	message any,
) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.InfoContext(ctx, "writing response")
	err := WriteJSON(w, status, ErrorMessage{Message: message}, nil)
	if err != nil {
		logger.ErrorContext(ctx, "error writing response", "error", err)
		logError(r, err)
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logError(r, err)
	message := "the server encountered a problem and could not process your request"
	logger.InfoContext(ctx, "the server encountered a problem and could not process your request")
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func WriteJSON(
	w http.ResponseWriter,
	status int,
	data any,
	headers http.Header,
) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	message := "the requested resource could not be found"
	logger.InfoContext(ctx, "the requested resource could not be found")
	errorResponse(w, r, http.StatusNotFound, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, message string) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.InfoContext(ctx, "bad request response")
	errorResponse(w, r, http.StatusBadRequest, message)
}
