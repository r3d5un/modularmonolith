package main

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/r3d5un/modularmonolith/internal/logging"
)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rCtx := r.Context()
		requestLogger := app.logger.With(
			slog.Group(
				"request",
				slog.String("id", uuid.New().String()),
				slog.String("method", r.Method),
				slog.String("protocol", r.Proto),
				slog.String("url", r.URL.Path),
			),
		)
		loggerCtx := logging.WithLogger(rCtx, requestLogger)
		requestLogger.Info("received request")

		next.ServeHTTP(w, r.WithContext(loggerCtx))
	})
}
