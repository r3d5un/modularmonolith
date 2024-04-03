package rabbit

import (
	"net/http"

	"github.com/r3d5un/modularmonolith/internal/logging"
	"github.com/r3d5un/modularmonolith/internal/utils/httputils"
)

type HelloWorldResponse struct {
	Message string `json:"message"`
}

func (m *Module) postHelloWorldMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("publishing message")
	err := m.queues.HelloWorldQueue.PublishHelloWorld()
	if err != nil {
		logger.Error("unable to publish message", "error", err)
		httputils.ServerErrorResponse(w, r, err)
		return
	}
	logger.Info("message published")

	logger.Info("writing response")
	err = httputils.WriteJSON(
		w,
		http.StatusOK,
		HelloWorldResponse{Message: "'Hello, World!' posted"},
		nil,
	)
	if err != nil {
		logger.Error("unable to write response", "error", err)
		httputils.ServerErrorResponse(w, r, err)
		return
	}
}
