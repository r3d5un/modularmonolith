package main

import (
	"net/http"

	"github.com/r3d5un/modularmonolith/internal/logging"
)

type ErrorMessage struct {
	Message any `json:"message"`
}

func (app *application) logError(r *http.Request, err error) {
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

func (app *application) errorResponse(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	message any,
) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.InfoContext(ctx, "writing response")
	err := app.writeJSON(w, status, ErrorMessage{Message: message}, nil)
	if err != nil {
		logger.ErrorContext(ctx, "error writing response", "error", err)
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	logger.InfoContext(ctx, "the server encountered a problem and could not process your request")
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}
