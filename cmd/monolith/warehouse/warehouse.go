package warehouse

import (
	"context"

	"github.com/r3d5un/modularmonolith/internal/monolith"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, mono monolith.Monolith) {}
