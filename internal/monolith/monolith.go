package monolith

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/r3d5un/modularmonolith/internal/config"
	"github.com/r3d5un/modularmonolith/internal/queue"
)

type Monolith interface {
	Config() *config.Configuration
	Logger() *slog.Logger
	DB() *sql.DB
	MQ() *queue.ChannelPool
	Mux() *http.ServeMux
}

type Module interface {
	Startup(ctx context.Context, mono Monolith)
}

type Warehouse interface{}
