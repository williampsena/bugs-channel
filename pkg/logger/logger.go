package logger

import (
	logrus "github.com/sirupsen/logrus"
	"github.com/williampsena/bugs-channel/pkg/config"
)

func Setup() {
	setupLogLevel()

	if config.IsProduction() {
		setupLogProduction()
	} else {
		setupLogDevelopment()
	}

}

func setupLogDevelopment() {
	logrus.SetFormatter(&logrus.TextFormatter{})
}

func setupLogProduction() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func setupLogLevel() {
	level, err := logrus.ParseLevel(config.LogLevel())

	if err != nil {
		logrus.Warnf("The log level %v is unsupported only: (trace, debug, info, warn, error, fatal, panic).", level)
		return
	}

	logrus.SetLevel(level)
}
