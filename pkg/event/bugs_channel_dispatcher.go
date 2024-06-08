// This package includes events contracts and behaviors
package event

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/williampsena/bugs-channel-plugins/pkg/event"
	"github.com/williampsena/bugs-channel/pkg/storage"
)

// A event of Dispatcher
type BugsChannelEventsDispatcher struct {
	queue storage.Queue
}

// Dispatch a event
func (d *BugsChannelEventsDispatcher) Dispatch(event event.Event) error {
	body, err := event.Json()

	if err != nil {
		return err
	}

	err = d.queue.Publish(context.TODO(), "events", body)

	if err != nil {
		return err
	}

	log.Infof("🐞 Ingest Event: %v", event.ID)

	return nil
}

// Dispatch many events to stdout
func (d *BugsChannelEventsDispatcher) DispatchMany(events []event.Event) error {
	for _, e := range events {
		err := d.Dispatch(e)

		if err != nil {
			return err
		}
	}

	return nil
}

// Creates a new event dispatcher
func NewDispatcher(queue storage.Queue) *BugsChannelEventsDispatcher {
	return &BugsChannelEventsDispatcher{queue}
}
