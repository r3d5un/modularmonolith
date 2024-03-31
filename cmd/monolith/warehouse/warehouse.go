package warehouse

import (
	"context"
	"log/slog"

	"github.com/r3d5un/modularmonolith/internal/monolith"
)

type Module struct {
	log *slog.Logger
}

func (m *Module) Startup(ctx context.Context, mono monolith.Monolith) {
	m.initModuleLogger(mono.Logger())
	m.log.Info("module logger initialized")

	// TODO: Add models
	// TODO: Add endpoints
}

func (m *Module) initModuleLogger(monoLogger *slog.Logger) {
	m.log = monoLogger.With(slog.Group("module", slog.String("name", "warehouse")))
}
