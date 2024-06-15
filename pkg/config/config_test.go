package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogLevel(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")

	require.Equal(t, LogLevel(), "debug")
}

func TestEnvironment(t *testing.T) {
	t.Setenv("GO_ENV", "development")

	require.Equal(t, Env(), "development")
}

func TestIsProduction(t *testing.T) {
	t.Setenv("GO_ENV", "production")
	require.Equal(t, IsProduction(), true)

	t.Setenv("GO_ENV", "test")
	require.Equal(t, IsProduction(), false)
}

func TestApiPort(t *testing.T) {
	t.Setenv("PORT", "1000")
	require.Equal(t, ApiPort(), 1000)
}

func TestConfigFille(t *testing.T) {
	t.Setenv("CONFIG_FILE", "/tmp/config.yml")
	require.Equal(t, ConfigFile(), "/tmp/config.yml")
}

func TestRateLimit(t *testing.T) {
	t.Setenv("WEB_RATE_LIMIT", "1")
	require.Equal(t, RateLimit(), int64(1))
}

func TestNatsConnectionUrl(t *testing.T) {
	t.Setenv("NATS_URL", "nats://localhost")
	require.Equal(t, NatsConnectionUrl(), "nats://localhost")
}

func TestRedisConnectionUrl(t *testing.T) {
	t.Setenv("REDIS_URL", "redis://localhost")
	require.Equal(t, RedisConnectionUrl(), "redis://localhost")
}

func TestEventChannel(t *testing.T) {
	t.Setenv("EVENT_CHANNEL", "nats_or_redis")
	require.Equal(t, EventChannel(), "nats_or_redis")
}
func TestScrubSensitiveKeys(t *testing.T) {
	t.Setenv("SCRUB_SENSITIVE_KEYS", "foo,bar")
	require.Equal(t, ScrubSensitiveKeys(), []string{"foo", "bar"})
}
