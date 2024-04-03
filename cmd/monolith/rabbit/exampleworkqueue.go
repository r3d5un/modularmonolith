package rabbit

import (
	"encoding/json"
	"net/http"

	"github.com/r3d5un/modularmonolith/internal/logging"
	"github.com/r3d5un/modularmonolith/internal/utils/httputils"
)

type ExampleWorkQueuePostBody struct {
	Message string `json:"message"`
}

type ExampleWorkQueueResponse struct {
	Message string `json:"message"`
}

func (m *Module) postExampleWorkQueueHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("unmarshalling message")
	var b ExampleWorkQueuePostBody
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		logger.Error("unable to decode HTTP request body", "body", r.Body)
		httputils.ServerErrorResponse(w, r, err)
		return
	}
	logger.Info("decoded request body", "message", "b")

	logger.Info("publishing message")
	err = m.queues.ExampleWorkQueue.PublishMessage(b.Message)
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
		ExampleWorkQueueResponse{Message: "message published"},
		nil,
	)
}
