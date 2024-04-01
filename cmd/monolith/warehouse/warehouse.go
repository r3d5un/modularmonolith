package warehouse

import (
	"context"
	"log/slog"
	"net/http"
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

	m.log.Info("registering routes")
	m.registerEndpoints(mono.Mux())
}

func (m *Module) initModuleLogger(monoLogger *slog.Logger) {
	m.log = monoLogger.With(slog.Group("module", slog.String("name", "warehouse")))
}

type RouteDefinition struct {
	Path    string
	Handler http.HandlerFunc
}

type RouteDefinitionList []RouteDefinition

func (m *Module) registerEndpoints(mux *http.ServeMux) {
	routes := RouteDefinitionList{
		{"GET /api/v1/warehouse/peppolbusinesscards", m.getPeppolBusinessCardHandler},
	}

	for _, d := range routes {
		m.log.Info("adding route", "route", d.Path)
		mux.Handle(d.Path, d.Handler)
	}
}
