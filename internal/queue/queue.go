package queue

import (
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ChannelPool struct {
	conn           *amqp.Connection
	channels       chan *amqp.Channel
	maxChannels    int
	shutdownSignal chan struct{}
}

func NewChannelPool(conn *amqp.Connection, maxChannels int) (*ChannelPool, error) {
	pool := &ChannelPool{
		conn:           conn,
		channels:       make(chan *amqp.Channel, maxChannels),
		maxChannels:    maxChannels,
		shutdownSignal: make(chan struct{}),
	}

	for i := 0; i < maxChannels; i++ {
		channel, err := conn.Channel()
		if err != nil {
			slog.Error("unable to create message queue channel", "error", err)
			return nil, err
		}

		pool.channels <- channel
	}

	return pool, nil
}

func (p *ChannelPool) GetChannel() (*amqp.Channel, error) {
	select {
	case channel := <-p.channels:
		return channel, nil
	case <-time.After(time.Second * 10):
		return nil, fmt.Errorf("unable to acquire channel")
	case <-p.shutdownSignal:
		return nil, fmt.Errorf("shutdown in progress")
	}
}

func (p *ChannelPool) ReturnChannel(channel *amqp.Channel) {
	select {
	case p.channels <- channel:
		// Channel returned to pool
	default:
		// Channel closed if pool is full
		channel.Close()
	}

}

func (p *ChannelPool) Shutdown() {
	close(p.shutdownSignal) // Signal all operations to shut down
	for len(p.channels) > 0 {
		channel := <-p.channels
		channel.Close() // Close all channels in the pool
	}
}
