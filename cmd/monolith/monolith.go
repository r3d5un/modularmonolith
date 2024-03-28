package main

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/r3d5un/modularmonolith/internal/config"
	"github.com/r3d5un/modularmonolith/internal/queue"
)

type application struct {
	cfg    *config.Configuration
	db     *sql.DB
	mq     *queue.ChannelPool
	mux    *http.ServeMux
	logger *slog.Logger
}

func (app *application) DB() *sql.DB {
	return app.db
}

func (app *application) MQ() *queue.ChannelPool {
	return app.mq
}

func (app *application) Mux() *http.ServeMux {
	return app.mux
}

func (app *application) Logger() *slog.Logger {
	return app.logger
}

func (app *application) Config() *config.Configuration {
	return app.cfg
}