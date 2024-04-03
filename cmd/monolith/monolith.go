package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/r3d5un/modularmonolith/internal/config"
	"github.com/r3d5un/modularmonolith/internal/monolith"
	"github.com/r3d5un/modularmonolith/internal/queue"
)

type application struct {
	cfg     *config.Configuration
	db      *sql.DB
	mq      *queue.ChannelPool
	mux     *http.ServeMux
	logger  *slog.Logger
	modules *monolith.Modules
	done    <-chan os.Signal
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

func (app *application) Modules() *monolith.Modules {
	return app.modules
}

func (app *application) Done() <-chan os.Signal {
	return app.done
}
