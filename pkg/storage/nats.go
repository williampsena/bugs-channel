package storage

import (
	"context"
	"errors"

	"github.com/nats-io/nats.go"
)

// Represents a Nats connection error
var ErrNatsConnection = errors.New("an error occurred while attempting to establish a Nats connection")

// Represents a Nats subscribe channel error
var ErrNatsSubscribeChannel = errors.New("an error occurred while attempting to subscribe to a Nats channel")

// The Nats instance
type Nats struct {
	conn *nats.Conn
}

// Build a new Nats connection
func NewNatsConnection(url string) (Queue, error) {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		return nil, errors.Join(ErrNatsConnection, err)
	}

	return &Nats{nc}, nil
}

// Publish a message
func (n *Nats) Publish(_ context.Context, topic string, message string) error {
	return n.conn.Publish(topic, []byte(message))
}

// Subscribe to a Nats channel
func (n *Nats) Subscribe(_ context.Context, channel string, handler SubscribeHandler) error {
	ch := make(chan *nats.Msg, 64)
	sub, err := n.conn.ChanSubscribe(channel, ch)

	if err != nil {
		return errors.Join(ErrNatsSubscribeChannel, err)
	}

	msg := <-ch

	handler(msg.Header, msg.Reply)

	sub.Unsubscribe()

	return nil
}

// Close Nats connection
func (n *Nats) Close() {
	n.conn.Drain()
	n.conn.Close()
}
