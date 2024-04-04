package main

import (
	"fmt"

	"github.com/r3d5un/modularmonolith/internal/config"
	"github.com/r3d5un/modularmonolith/internal/queue"
	amqp "github.com/rabbitmq/amqp091-go"
)

func openQueue(config config.MessageQueueConfiguration) (mqPool *queue.ChannelPool, err error) {
	conn, err := amqp.Dial(config.DSN)
	fmt.Println(config.DSN)
	if err != nil {
		return nil, err
	}

	mqPool, err = queue.NewChannelPool(conn, config.MaxConns)
	if err != nil {
		return nil, err
	}

	return mqPool, nil
}
