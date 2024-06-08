package storage

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/williampsena/bugs-channel/pkg/config"
)

// Represents a Redis connection error
var ErrRedisConnection = errors.New("an error occurred while attempting to establish a Redis connection")

// Represents a Redis subscribe channel error
var ErrRedisSubscribeChannel = errors.New("an error occurred while attempting to subscribe to a Redis channel")

// The Redis instance
type Redis struct {
	conn *redis.Client
}

// Build a new Redis connection
func NewRedisConnection(url string) (Queue, error) {
	opts, err := buildRedisOptions(config.RedisConnectionUrl())

	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opts)

	return &Redis{rdb}, nil
}

// Build redis connection options
func buildRedisOptions(url string) (*redis.Options, error) {
	opts, err := redis.ParseURL(url)

	if err != nil {
		return nil, errors.Join(ErrRedisConnection, err)
	}

	return opts, nil
}

// Publish a message
func (r *Redis) Publish(ctx context.Context, topic string, message string) error {
	return r.conn.Publish(ctx, topic, message).Err()
}

// Subscribe to a Redis topic
func (r *Redis) Subscribe(ctx context.Context, channel string, handler SubscribeHandler) error {
	pubsub := r.conn.Subscribe(ctx, channel)

	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		handler(buildRedisHeaders(msg.Channel, msg.Pattern), msg.Payload)
	}

	return nil
}

func buildRedisHeaders(channel string, pattern string) map[string][]string {
	return map[string][]string{
		"channel": {channel},
		"pattern": {pattern},
	}
}

// Close Redis connection
func (r *Redis) Close() {
	r.Close()
}
