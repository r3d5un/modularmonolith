package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	app.logger.Info("creating standard middleware chain")
	standard := alice.New(app.recoverPanic, app.logRequest)

	handler := standard.Then(app.mux)
	return handler
}
