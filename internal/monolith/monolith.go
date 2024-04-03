package monolith

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/r3d5un/modularmonolith/internal/config"
	"github.com/r3d5un/modularmonolith/internal/peppol"
	"github.com/r3d5un/modularmonolith/internal/queue"
)

type Monolith interface {
	Config() *config.Configuration
	Logger() *slog.Logger
	DB() *sql.DB
	MQ() *queue.ChannelPool
	Mux() *http.ServeMux
	Modules() *Modules
	Done() <-chan os.Signal
}

type Modules struct {
	Rabbit    Rabbit
	Warehouse Warehouse
}

type Module interface {
	// Setup sets up the module using the context and resources from the
	// monolith. For example, initializing database models, message queues
	// and so on.
	//
	// Be aware that in case of injecting resources from another module,
	// no method call should be made within the Setup process. This can
	// cause segmentation faults as there is no guarantee that injected modules
	// have completed their own startup process.
	//
	// If resources are required from other modules as part of the startup
	// process, add a PostSetup or Startup method, and set it to run after
	// Setup has completed.
	Setup(ctx context.Context, mono Monolith)
	// PostSetup performs any additional setup tasks, usually in cases where
	// external resources from other modules were required, and a guarantee
	// that the initial setup of any given module is completed to avoid
	// segmentation faults.
	PostSetup()
}

type Warehouse interface {
	GetPeppolBusinessCard(ctx context.Context, id string) (*peppol.BusinessCard, error)
}

type Rabbit interface{}
