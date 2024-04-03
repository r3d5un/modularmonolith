package queue

import (
	"log/slog"

	"github.com/r3d5un/modularmonolith/internal/queue"
)

type Queues struct {
	HelloWorldQueue HelloWorldQueue
}

func NewQueues(pool *queue.ChannelPool) (*Queues, error) {
	helloWorldQueue, err := NewHelloWorldQueue(pool)
	if err != nil {
		slog.Error("unable to create a new queue", "error", err)
		return nil, err
	}

	qs := Queues{
		HelloWorldQueue: *helloWorldQueue,
	}

	return &qs, nil
}
