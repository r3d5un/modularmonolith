package queue

import (
	"log/slog"
	"os"

	"github.com/r3d5un/modularmonolith/internal/queue"
)

type Queues struct {
	HelloWorldQueue  HelloWorldQueue
	ExampleWorkQueue ExampleWorkQueue
}

func NewQueues(pool *queue.ChannelPool) (*Queues, error) {
	helloWorldQueue, err := NewHelloWorldQueue(pool)
	if err != nil {
		slog.Error("unable to create a new queue", "error", err)
		return nil, err
	}

	exampleWorkQueue, err := NewExampleWorkQueue(pool)
	if err != nil {
		slog.Error("unable to create new queue", "error", err)
		return nil, err
	}

	qs := Queues{
		HelloWorldQueue:  *helloWorldQueue,
		ExampleWorkQueue: *exampleWorkQueue,
	}

	return &qs, nil
}

func ConsumeExampleWorkQueue(queue ExampleWorkQueue, done <-chan os.Signal) {
	msgs, err := queue.GetMessages(false)
	if err != nil {
		slog.Error("Error getting messages", "error", err)
		return
	}

	for {
		select {
		case msg := <-msgs:
			slog.Info("Received message", "body", string(msg.Body))
			err := msg.Ack(false)
			if err != nil {
				slog.Error("Error acknowledging message", "error", err)
			}
		case <-done:
			slog.Info("received shutdown signal")
			return
		}
	}
}
