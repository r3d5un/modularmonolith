package rabbit

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/r3d5un/modularmonolith/internal/monolith"
	"github.com/r3d5un/modularmonolith/internal/rabbit/queues"
	"github.com/r3d5un/modularmonolith/internal/requests"
)

type Module struct {
	log      *slog.Logger
	wh       monolith.Warehouse
	queues   *queue.Queues
	whClient monolith.Warehouse
}

func (m *Module) Setup(ctx context.Context, mono monolith.Monolith) {
	m.initModuleLogger(mono.Logger())
	m.log.Info("module logger initialized")

	m.log.Info("injecting warehouse module")
	m.wh = mono.Modules().Warehouse

	m.log.Info("creating warehouse client")
	m.whClient = requests.NewClient(mono.Config().Rq)

	m.log.Info("initializing queues")
	queues, err := queue.NewQueues(mono.MQ())
	if err != nil {
		m.log.Error("unable to initialize new queues", "error", err)
		os.Exit(1)
	}
	m.queues = queues

	m.log.Info("creating background processes")
	go queue.ConsumeExampleWorkQueue(m.queues.ExampleWorkQueue, mono.Done())

	m.log.Info("registering routes")
	m.registerEndpoints(mono.Mux())
}

func (m *Module) PostSetup() {
	m.log.Info("performing post setup process")

	m.log.Info("making example function call to warehouse")
	pbc, err := m.wh.GetPeppolBusinessCard(context.Background(), "0088:5903351900034")
	if err != nil {
		m.log.Error("unable to get peppol business card", "error", err)
		return
	}
	m.log.Info("retrieved peppol business card", "pbc", pbc)
}

func (m *Module) initModuleLogger(monoLogger *slog.Logger) {
	m.log = monoLogger.With(slog.Group("module", slog.String("name", "rabbit")))
}

type RouteDefinition struct {
	Path    string
	Handler http.HandlerFunc
}

type RouteDefinitionList []RouteDefinition

func (m *Module) registerEndpoints(mux *http.ServeMux) {
	routes := RouteDefinitionList{
		{"POST /api/v1/queue/hello_world", m.postHelloWorldMessageHandler},
		{"POST /api/v1/queue/example_work_queue", m.postExampleWorkQueueHandler},
	}

	for _, d := range routes {
		m.log.Info("adding route", "route", d.Path)
		mux.Handle(d.Path, d.Handler)
	}
}
