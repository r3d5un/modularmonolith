package rabbit

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/r3d5un/modularmonolith/internal/monolith"
)

type Module struct {
	log *slog.Logger
	wh  monolith.Warehouse
}

func (m *Module) Setup(ctx context.Context, mono monolith.Monolith) {
	m.initModuleLogger(mono.Logger())
	m.log.Info("module logger initialized")

	m.log.Info("injecting warehouse module")
	m.wh = mono.Modules().Warehouse
}

func (m *Module) PostSetup() {
	m.log.Info("performing post setup process")
}

func (m *Module) initModuleLogger(monoLogger *slog.Logger) {
	m.log = monoLogger.With(slog.Group("module", slog.String("name", "rabbit")))
}

type RouteDefinition struct {
	Path    string
	Handler http.HandlerFunc
}

type RouteDefinitionList []RouteDefinition

func (m *Module) registerEndpoints(mux *http.ServeMux) {}
