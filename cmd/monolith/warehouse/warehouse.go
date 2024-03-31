package warehouse

import (
	"context"
	"log/slog"
	"time"

	"github.com/r3d5un/modularmonolith/internal/monolith"
	"github.com/r3d5un/modularmonolith/internal/warehouse/data"
)

type Module struct {
	log    *slog.Logger
	models data.Models
}

func (m *Module) Startup(ctx context.Context, mono monolith.Monolith) {
	m.initModuleLogger(mono.Logger())
	m.log.Info("module logger initialized")

	m.log.Info("adding warehouse data models")
	timeout := time.Duration(mono.Config().DB.Timeout) * time.Second
	m.models = data.NewModels(mono.DB(), &timeout)

	// TODO: Add endpoints
}

func (m *Module) initModuleLogger(monoLogger *slog.Logger) {
	m.log = monoLogger.With(slog.Group("module", slog.String("name", "warehouse")))
}
