// This package contains Stprage implementations such as Nats
package storage

import "context"

// The queue interface
type Queue interface {
	// Publish a message
	Publish(ctx context.Context, topic string, message string) error

	// Subscribe to a topic
	Subscribe(ctx context.Context, topic string, handler SubscribeHandler) error
}

// The subscribe handler signature
type SubscribeHandler func(header map[string][]string, body string) error
