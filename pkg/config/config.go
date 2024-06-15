// This package provides environment variable reading support.
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// This occurs when an invalid server port number is provided.
var ErrInvalidPort = errors.New("the port is invalid")

const (
	// Production environment enumerator
	production = "production"
)

// Returns the current log level
func LogLevel() string {
	return os.Getenv("LOG_LEVEL")
}

// Returns the current environment
func Env() string {
	return os.Getenv("GO_ENV")
}

// Define if the current environment is production
func IsProduction() bool {
	return Env() == production
}

// Returns the listen api port application
func ApiPort() int {
	env := getEnv("PORT", "4000")

	port, err := strconv.Atoi(env)

	if err != nil {
		panic(errors.Join(ErrInvalidPort, err))
	}

	return port
}

// Returns the config file path
func ConfigFile() string {
	return os.Getenv("CONFIG_FILE")
}

// Requests rate limit
func RateLimit() int64 {
	value, err := strconv.ParseInt(getEnv("WEB_RATE_LIMIT", "0"), 10, 64)

	if err != nil {
		return 0
	}

	return value
}

// The Nats connection url
func NatsConnectionUrl() string {
	return os.Getenv("NATS_URL")
}

// The Redis connection url
func RedisConnectionUrl() string {
	return os.Getenv("REDIS_URL")
}

// The event channel (redis/nats)
func EventChannel() string {
	return os.Getenv("EVENT_CHANNEL")
}

// The sensitive keys to hide from events
func ScrubSensitiveKeys() []string {
	return strings.Split(getEnv("SCRUB_SENSITIVE_KEYS", ""), ",")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
