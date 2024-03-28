package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	app.logger.Info("creating standard middleware chain")
	// TODO: middleware logRequest
	// TODO: middleware recoverPanic
	standard := alice.New()

	handler := standard.Then(app.mux)
	return handler
}
