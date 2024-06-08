package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/williampsena/bugs-channel-plugins/pkg/sentry"
	"github.com/williampsena/bugs-channel/pkg/config"
	"github.com/williampsena/bugs-channel/pkg/event"
	"github.com/williampsena/bugs-channel/pkg/logger"
	"github.com/williampsena/bugs-channel/pkg/service"
	"github.com/williampsena/bugs-channel/pkg/settings"
	"github.com/williampsena/bugs-channel/pkg/storage"
	"github.com/williampsena/bugs-channel/pkg/web"
)

func init() {
	logger.Setup()
}

func main() {
	configFile, err := settings.BuildConfigFile(config.ConfigFile())

	if err != nil {
		log.Fatal("❌ The configuration file is in incorrect format or does not exist.", err)
	}

	nats := buildQueue()

	sentryServerContext := sentry.ServerContext{
		Context:          context.Background(),
		ServiceFetcher:   service.NewYAMLServiceFetcher(configFile.Services),
		EventsDispatcher: event.NewDispatcher(nats),
	}

	sentrySvr := sentry.BuildServer(&sentryServerContext)
	go sentry.SetupServer(sentrySvr)

	webServerContext := web.ServerContext{
		Context: context.Background(),
		Queue:   nats,
	}

	web.SetupServer(&webServerContext)
}

func buildQueue() storage.Queue {
	var queue storage.Queue
	var err error

	eventChannel := config.EventChannel()

	log.Infof("The event channel is %v", eventChannel)

	switch eventChannel {
	case "nats":
		queue, err = storage.NewNatsConnection(config.NatsConnectionUrl())
	case "redis":
		queue, err = storage.NewRedisConnection(config.RedisConnectionUrl())
	}

	if err != nil {
		log.Fatal("❌ Something went wrong when trying to construct Queue's connection.", err)
	}

	return queue
}
