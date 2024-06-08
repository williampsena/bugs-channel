package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/williampsena/bugs-channel-plugins/pkg/event"
	"github.com/williampsena/bugs-channel-plugins/pkg/test"
	"github.com/williampsena/bugs-channel/pkg/storage"
)

func TestDispatchSuccess(t *testing.T) {
	buf := test.CaptureLog()

	dispatcher := NewDispatcher(buildTestQueue())

	err := dispatcher.Dispatch(event.Event{
		ID:        "foo",
		ServiceId: "bar",
		Platform:  "python",
	})

	require.Nil(t, err)

	assert.Contains(t, buf.String(), "🐞 Ingest Event: foo")

	test.ResetCaptureLog()
}

type mockNats struct {
	lastMessage string
}

func buildTestQueue() storage.Queue {
	return &mockNats{lastMessage: ""}
}

// Publish a message
func (n *mockNats) Publish(ctx context.Context, topic string, message string) error {
	n.lastMessage = message
	return nil
}

// Subscribe to a Nats channel
func (n *mockNats) Subscribe(ctx context.Context, channel string, handler storage.SubscribeHandler) error {
	var header = map[string][]string{"foo": []string{"bar"}}
	handler(header, "foo")
	return nil
}

// Close Nats connection
func (n *mockNats) Close() {
	n.lastMessage = ""
}
