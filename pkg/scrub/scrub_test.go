package scrub

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/williampsena/bugs-channel-plugins/pkg/event"
)

func TestScrubSensitiveEventSimple(t *testing.T) {
	simpleEvent := event.Event{
		ID:        "foo",
		ServiceId: "bar",
		Platform:  "python",
		Extra:     event.EventExtra{"message": "public info", "password": "foo", "secret": "bar"},
		StackTrace: event.StackTrace{
			map[string]interface{}{"error": "fatal", "secret": "public"},
		},
		Tags: []string{"secret:biz", "password:qux", "app:foo"},
	}

	ScrubSensitiveEvent(&simpleEvent, []string{"password", "secret"})

	expectedEvent := event.Event{
		ID:        "foo",
		ServiceId: "bar",
		Platform:  "python",
		Extra:     event.EventExtra{"message": "public info", "password": "*", "secret": "*"},
		StackTrace: event.StackTrace{
			map[string]interface{}{"error": "fatal", "secret": "public"},
		},
		Tags: []string{"secret:*", "password:*", "app:foo"},
	}

	assert.Equal(t, expectedEvent, simpleEvent)
}

func TestScrubSensitiveEventComplex(t *testing.T) {
	complexEvent := event.Event{
		ID:        "foo",
		ServiceId: "bar",
		Platform:  "python",
		Extra: event.EventExtra{
			"message": "public info",
			"user":    map[string]interface{}{"credentials": map[string]interface{}{"pwd": "foo"}},
			"keys": []map[string]interface{}{
				{"secret": "bar"},
			},
		},
		StackTrace: event.StackTrace{
			map[string]interface{}{"error": "fatal", "secret": "public"},
		},
		Tags: []string{"secret:biz", "pwd:qux", "app:foo"},
	}

	ScrubSensitiveEvent(&complexEvent, []string{"pwd", "secret"})

	expectedEvent := event.Event{
		ID:        "foo",
		ServiceId: "bar",
		Platform:  "python",
		Extra: event.EventExtra{
			"message": "public info",
			"user":    map[string]interface{}{"credentials": map[string]interface{}{"pwd": "*"}},
			"keys":    []map[string]interface{}{{"secret": "*"}},
		},
		StackTrace: event.StackTrace{
			map[string]interface{}{"error": "fatal", "secret": "public"},
		},
		Tags: []string{"secret:*", "pwd:*", "app:foo"},
	}

	assert.Equal(t, expectedEvent, complexEvent)
}
